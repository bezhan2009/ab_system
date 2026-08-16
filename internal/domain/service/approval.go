package service

import (
	"ab_system/internal/clients/telegram_notifications"
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/errs"
	"context"
	"fmt"
)

type ApprovalService struct {
	approvalReader  repository.ApprovalReader
	approvalWriter  repository.ApprovalWriter
	approvalDeleter repository.ApprovalDeleter

	experimentReader repository.ExperimentReader
	experimentWriter repository.ExperimentWriter

	approverGroupReader repository.ApproverGroupReader
	approverGroupWriter repository.ApproverGroupWriter

	userReader repository.UserReader

	telegramClient *telegram_notifications.NotifyClient
}

func NewApprovalService(
	approvalReader repository.ApprovalReader,
	approvalWriter repository.ApprovalWriter,
	approvalDeleter repository.ApprovalDeleter,
	experimentReader repository.ExperimentReader,
	experimentWriter repository.ExperimentWriter,
	approverGroupReader repository.ApproverGroupReader,
	approverGroupWriter repository.ApproverGroupWriter,
	userReader repository.UserReader,
	telegramClient *telegram_notifications.NotifyClient,
) *ApprovalService {
	return &ApprovalService{
		approvalReader:      approvalReader,
		approvalWriter:      approvalWriter,
		approvalDeleter:     approvalDeleter,
		experimentReader:    experimentReader,
		experimentWriter:    experimentWriter,
		approverGroupReader: approverGroupReader,
		approverGroupWriter: approverGroupWriter,
		userReader:          userReader,
		telegramClient:      telegramClient,
	}
}

func (s *ApprovalService) GetApprovalsByExperimentID(ctx context.Context, experimentID string) (approvals *[]models.Approval, err error) {
	approvals, err = s.approvalReader.GetApprovalsByExperimentID(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	return approvals, nil
}

func (s *ApprovalService) GetApprovalByID(ctx context.Context, approvalID string) (approval *models.Approval, err error) {
	approval, err = s.approvalReader.GetApprovalByID(ctx, approvalID)
	if err != nil {
		return nil, err
	}

	return approval, nil
}

func (s *ApprovalService) CreateApproval(ctx context.Context, approval *models.Approval) (err error) {
	exp, err := s.experimentReader.GetExperimentByID(ctx, approval.ExperimentID.String())
	if err != nil {
		return err
	}

	if exp.Status != models.StatusInReview && exp.Status != models.StatusApproved {
		return errs.ErrExperimentStatusIsNotOnReview
	}

	approverGroup, members, err := s.approverGroupReader.GetApproverGroupByExperimentID(ctx, approval.ExperimentID.String())
	if err != nil {
		return err
	}

	isApprover := false
	for _, member := range members {
		if member.ID == approval.ApproverID {
			isApprover = true
			break
		}
	}

	if !isApprover {
		return errs.ErrForbiddenToBeApprover
	}

	existingApprovals, err := s.GetApprovalsByExperimentID(ctx, approval.ExperimentID.String())
	if err != nil {
		return err
	}
	for _, a := range *existingApprovals {
		if a.ApproverID == approval.ApproverID && exp.Version == a.ExperimentVersion {
			return errs.ErrApproverAlreadyVoted
		}
	}

	approval.ExperimentVersion = exp.Version

	err = s.approvalWriter.CreateApproval(ctx, approval)
	if err != nil {
		return err
	}

	if !approval.Approved {
		exp.Status = models.StatusRejected

		_, err = s.experimentWriter.UpdateExperiment(ctx, exp)
		if err != nil {
			return err
		}
	} else {
		approvalCount := 0
		for _, a := range *existingApprovals {
			if a.Approved {
				approvalCount++
			}
		}

		approvalCount++

		if approvalCount >= approverGroup.MinApprovals {
			exp.Status = models.StatusApproved
			_, err = s.experimentWriter.UpdateExperiment(ctx, exp)
			if err != nil {
				return err
			}
		}
	}

	go func() {
		bgCtx := context.Background()

		ownerEmail := ""

		owner, err := s.userReader.GetUserByID(bgCtx, exp.OwnerID)
		if err != nil {
			ownerEmail = exp.OwnerID
		} else {
			ownerEmail = fmt.Sprintf("%s (%s)", owner.Email, exp.OwnerID)
		}

		_ = SendTelegramNotificationByExperimentStatus(bgCtx, s.telegramClient.Notify, exp, ownerEmail)
	}()

	return nil
}

func (s *ApprovalService) DeleteApproval(ctx context.Context, approvalID string) (err error) {
	err = s.approvalDeleter.DeleteApproval(ctx, approvalID)
	if err != nil {
		return err
	}

	return nil
}
