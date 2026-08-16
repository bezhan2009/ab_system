package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApproverGroupRepository struct {
	db *gorm.DB
}

func NewApproverGroupRepository(db *gorm.DB) *ApproverGroupRepository {
	return &ApproverGroupRepository{db: db}
}

func (r *ApproverGroupRepository) GetApproverGroupByExperimenterID(ctx context.Context, experimenterID string) (approverGroup *models.ApproverGroup, err error) {
	const op = "repository.postgres.GetApproverGroupByExperimenterID"

	uid, err := uuid.Parse(experimenterID)
	if err != nil {
		return nil, err
	}

	if err = r.db.WithContext(ctx).Where("experimenter_id = ?", uid).First(&approverGroup).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return approverGroup, nil
}

func (r *ApproverGroupRepository) GetApproverGroupWithMembers(ctx context.Context, groupID string) (approverGroup *models.ApproverGroup, members []models.User, err error) {
	const op = "repository.postgres.GetApproverGroupWithMembers"

	gid, err := uuid.Parse(groupID)
	if err != nil {
		return nil, nil, err
	}

	if err = r.db.WithContext(ctx).Where("id = ?", gid).First(&approverGroup).Error; err != nil {
		return nil, nil, TranslateGormError(err)
	}

	err = r.db.WithContext(ctx).
		Joins("JOIN approver_group_members ON approver_group_members.approver_id = users.id").
		Where("approver_group_members.approver_group_id = ?", gid).
		Find(&members).Error
	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error loading members: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, nil, TranslateGormError(err)
	}

	return approverGroup, members, nil
}

func (r *ApproverGroupRepository) GetApproverGroupByExperimentID(ctx context.Context, experimentID string) (approverGroup *models.ApproverGroup, members []models.User, err error) {
	const op = "repository.postgres.GetApproverGroupByExperimentID"

	expID, err := uuid.Parse(experimentID)
	if err != nil {
		return nil, nil, err
	}

	var experiment models.Experiment
	if err = r.db.WithContext(ctx).
		Where("id = ?", expID).
		First(&experiment).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error getting experiment: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, nil, TranslateGormError(err)
	}

	ownerID, err := uuid.Parse(experiment.OwnerID)
	if err != nil {
		return nil, nil, err
	}

	if err = r.db.WithContext(ctx).
		Where("experimenter_id = ?", ownerID).
		First(&approverGroup).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error getting approver group: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, nil, TranslateGormError(err)
	}

	err = r.db.WithContext(ctx).
		Joins("JOIN approver_group_members ON approver_group_members.approver_id = users.id").
		Where("approver_group_members.approver_group_id = ?", approverGroup.ID).
		Find(&members).Error

	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error loading members: %v",
			op, observability.GetTraceID(ctx), err)
		return nil, nil, TranslateGormError(err)
	}

	return approverGroup, members, nil
}

func (r *ApproverGroupRepository) CreateApproverGroup(ctx context.Context, group *models.ApproverGroup, approverIDs []string) (err error) {
	const op = "repository.postgres.CreateApproverGroup"

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Create(group).Error; err != nil {
			logger.Error.Printf("[%s] TraceId=%s Error approver group: %v",
				op, observability.GetTraceID(ctx), err)

			return err
		}

		for _, aid := range approverIDs {
			uid, err := uuid.Parse(aid)
			if err != nil {
				return err
			}

			member := models.ApproverGroupMember{
				ApproverGroupID: group.ID,
				ApproverID:      uid,
			}

			if err = tx.Create(&member).Error; err != nil {
				logger.Error.Printf("[%s] TraceId=%s Error creating members: %v",
					op, observability.GetTraceID(ctx), err)

				return err
			}
		}

		return nil
	})
}

func (r *ApproverGroupRepository) UpdateApproverGroup(ctx context.Context, group *models.ApproverGroup, approverIDs []string) (approverGroup *models.ApproverGroup, err error) {
	const op = "repository.postgres.UpdateApproverGroup"
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Model(&models.ApproverGroup{}).Where("id = ?", group.ID).Updates(group).Error; err != nil {
			return err
		}

		if err = tx.Where("approver_group_id = ?", group.ID).Delete(&models.ApproverGroupMember{}).Error; err != nil {
			return err
		}

		for _, aid := range approverIDs {
			uid, err := uuid.Parse(aid)
			if err != nil {
				return err
			}

			member := models.ApproverGroupMember{
				ApproverGroupID: group.ID,
				ApproverID:      uid,
			}

			if err = tx.Create(&member).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v", op, observability.GetTraceID(ctx), err)
		return nil, TranslateGormError(err)
	}

	return group, nil
}

func (r *ApproverGroupRepository) DeleteApproverGroup(ctx context.Context, id string) (err error) {
	const op = "repository.postgres.DeleteApproverGroup"
	gid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Where("approver_group_id = ?", gid).Delete(&models.ApproverGroupMember{}).Error; err != nil {
			logger.Error.Printf("[%s] TraceId=%s Error deleting approver group: %v",
				op, observability.GetTraceID(ctx), err)

			return err
		}

		return tx.Delete(&models.ApproverGroup{}, "id = ?", gid).Error
	})
}
