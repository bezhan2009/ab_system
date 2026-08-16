package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/errs"
	"context"
)

type ApproverGroupService struct {
	groupReader repository.ApproverGroupReader
	groupWriter repository.ApproverGroupWriter
	userReader  repository.UserReader
}

func NewApproverGroupService(
	groupReader repository.ApproverGroupReader,
	groupWriter repository.ApproverGroupWriter,
	userReader repository.UserReader,
) *ApproverGroupService {
	return &ApproverGroupService{
		groupReader: groupReader,
		groupWriter: groupWriter,
		userReader:  userReader,
	}
}

func (s *ApproverGroupService) GetGroupByExperimenterID(ctx context.Context, experimenterID string) (*models.ApproverGroup, []models.User, error) {
	group, err := s.groupReader.GetApproverGroupByExperimenterID(ctx, experimenterID)
	if err != nil {
		return nil, nil, err
	}
	if group == nil {
		return nil, nil, nil
	}
	_, members, err := s.groupReader.GetApproverGroupWithMembers(ctx, group.ID.String())
	if err != nil {
		return nil, nil, err
	}
	return group, members, nil
}

func (s *ApproverGroupService) GetApproverGroupByExperimentID(ctx context.Context, experimentID string) (approverGroup *models.ApproverGroup, members []models.User, err error) {
	approverGroup, members, err = s.groupReader.GetApproverGroupByExperimentID(ctx, experimentID)
	if err != nil {
		return nil, nil, err
	}

	return approverGroup, members, nil
}

func (s *ApproverGroupService) CreateApproverGroup(ctx context.Context, group *models.ApproverGroup, approverIDs []string) error {
	exp, err := s.userReader.GetUserByID(ctx, group.ExperimenterID.String())
	if err != nil {
		return err
	}
	if exp.Role != "experimenter" && exp.Role != "admin" {
		return errs.ErrRoleMustBeExperimenterOrAdmin
	}

	for _, aid := range approverIDs {
		app, err := s.userReader.GetUserByID(ctx, aid)
		if err != nil {
			return err
		}
		if app.Role != "approver" && app.Role != "admin" {
			return errs.ErrInvalidRole
		}
	}

	existing, _ := s.groupReader.GetApproverGroupByExperimenterID(ctx, group.ExperimenterID.String())
	if existing != nil {
		return errs.ErrAlreadyExists
	}
	return s.groupWriter.CreateApproverGroup(ctx, group, approverIDs)
}

func (s *ApproverGroupService) UpdateApproverGroup(ctx context.Context, group *models.ApproverGroup, approverIDs []string) (*models.ApproverGroup, error) {
	existing, err := s.groupReader.GetApproverGroupByExperimenterID(ctx, group.ExperimenterID.String())
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errs.ErrRecordNotFound
	}
	group.ID = existing.ID

	exp, err := s.userReader.GetUserByID(ctx, group.ExperimenterID.String())
	if err != nil {
		return nil, err
	}

	if exp.Role != "experimenter" && exp.Role != "admin" {
		return nil, errs.ErrRoleMustBeExperimenterOrAdmin
	}

	for _, aid := range approverIDs {
		app, err := s.userReader.GetUserByID(ctx, aid)
		if err != nil {
			return nil, err
		}
		if app.Role != "approver" && app.Role != "admin" {
			return nil, errs.ErrInvalidRole
		}
	}

	return s.groupWriter.UpdateApproverGroup(ctx, group, approverIDs)
}

func (s *ApproverGroupService) DeleteApproverGroup(ctx context.Context, id string) error {
	return s.groupWriter.DeleteApproverGroup(ctx, id)
}
