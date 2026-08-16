package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/errs"
	"context"
	"errors"
)

type TeamService struct {
	teamReader  repository.TeamReader
	teamWriter  repository.TeamWriter
	teamDeleter repository.TeamDeleter

	userReader repository.UserReader
	userWriter repository.UserWriter
}

func NewTeamService(
	teamReader repository.TeamReader,
	teamWriter repository.TeamWriter,
	teamDeleter repository.TeamDeleter,
	userReader repository.UserReader,
	userWriter repository.UserWriter) *TeamService {
	return &TeamService{
		teamReader:  teamReader,
		teamWriter:  teamWriter,
		teamDeleter: teamDeleter,
		userReader:  userReader,
		userWriter:  userWriter,
	}
}

func (s *TeamService) GetAllTeams(ctx context.Context, name string) (*[]models.Team, error) {
	if name != "" {
		return s.teamReader.GetTeamByNameLike(ctx, name)
	} else {
		return s.teamReader.GetAllTeams(ctx)
	}
}

func (s *TeamService) GetTeamByID(ctx context.Context, id string) (*models.Team, error) {
	return s.teamReader.GetTeamByID(ctx, id)
}

func (s *TeamService) CreateTeam(ctx context.Context, team *models.Team) error {
	existing, err := s.teamReader.GetTeamByName(ctx, team.Name)
	if err != nil {
		if !errors.Is(err, errs.ErrRecordNotFound) {
			return err
		}
	}
	if existing != nil {
		return errs.ErrAlreadyExists
	}

	return s.teamWriter.CreateTeam(ctx, team)
}

func (s *TeamService) UpdateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	existing, err := s.teamReader.GetTeamByID(ctx, team.ID.String())
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errs.ErrRecordNotFound
	}

	if team.Name != existing.Name {
		dup, err := s.teamReader.GetTeamByName(ctx, team.Name)
		if err != nil {
			if !errors.Is(err, errs.ErrRecordNotFound) {
				return nil, err
			}
		}
		if dup != nil && dup.ID != team.ID {
			return nil, errs.ErrAlreadyExists
		}
	}

	return s.teamWriter.UpdateTeam(ctx, team)
}

func (s *TeamService) AddTeamMember(ctx context.Context, teamID string, userID string) (err error) {
	team, err := s.teamReader.GetTeamByID(ctx, teamID)
	if err != nil {
		return err
	}

	user, err := s.userReader.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.TeamID = &team.ID

	_, err = s.userWriter.UpdateUser(ctx, user)

	return err
}

func (s *TeamService) RemoveTeamMember(ctx context.Context, teamID string, userID string) (err error) {
	_, err = s.teamReader.GetTeamByID(ctx, teamID)
	if err != nil {
		return err
	}

	user, err := s.userReader.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.TeamID = nil

	_, err = s.userWriter.UpdateUser(ctx, user)
	return err
}

func (s *TeamService) DeleteTeam(ctx context.Context, teamId string) error {
	return s.teamDeleter.DeleteTeam(ctx, teamId)
}
