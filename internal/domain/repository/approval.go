package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type ApprovalReader interface {
	GetApprovalsByExperimentID(ctx context.Context, experimentID string) (approvals *[]models.Approval, err error)
	GetApprovalByID(ctx context.Context, approvalID string) (approval *models.Approval, err error)
}

type ApprovalWriter interface {
	CreateApproval(ctx context.Context, approval *models.Approval) (err error)
}

type ApprovalDeleter interface {
	DeleteApproval(ctx context.Context, approvalID string) (err error)
}
