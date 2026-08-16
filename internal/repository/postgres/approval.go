package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"

	"gorm.io/gorm"
)

type ApprovalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) *ApprovalRepository {
	return &ApprovalRepository{
		db: db,
	}
}

func (r *ApprovalRepository) GetApprovalsByExperimentID(ctx context.Context, experimentID string) (approvals *[]models.Approval, err error) {
	const op = "repository.postgres.GetApprovalsByExperimentID"

	if err = r.db.WithContext(ctx).Find(&approvals, "experiment_id = ?", experimentID).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return approvals, nil
}

func (r *ApprovalRepository) GetApprovalByID(ctx context.Context, approvalID string) (approval *models.Approval, err error) {
	const op = "repository.postgres.GetApprovalByID"

	if err = r.db.WithContext(ctx).First(&approval, "id = ?", approvalID).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return approval, nil
}

func (r *ApprovalRepository) CreateApproval(ctx context.Context, approval *models.Approval) (err error) {
	const op = "repository.postgres.CreateApproval"

	if err = r.db.WithContext(ctx).Create(&approval).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *ApprovalRepository) DeleteApproval(ctx context.Context, approvalID string) (err error) {
	const op = "repository.postgres.DeleteApproval"

	if err = r.db.WithContext(ctx).Where("id = ?", approvalID).Delete(&models.Approval{}).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}
