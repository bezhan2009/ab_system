package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type ApproverGroupReader interface {
	GetApproverGroupByExperimenterID(ctx context.Context, experimenterID string) (approverGroup *models.ApproverGroup, err error)
	GetApproverGroupWithMembers(ctx context.Context, groupID string) (approverGroup *models.ApproverGroup, members []models.User, err error)
	GetApproverGroupByExperimentID(ctx context.Context, experimentID string) (approverGroup *models.ApproverGroup, members []models.User, err error)
}

type ApproverGroupWriter interface {
	CreateApproverGroup(ctx context.Context, group *models.ApproverGroup, approverIDs []string) (err error)
	UpdateApproverGroup(ctx context.Context, group *models.ApproverGroup, approverIDs []string) (approverGroup *models.ApproverGroup, err error)
	DeleteApproverGroup(ctx context.Context, approverGroupId string) (err error)
}
