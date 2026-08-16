package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type TeamReader interface {
	GetTeamByID(ctx context.Context, teamId string) (team *models.Team, err error)
	GetAllTeams(ctx context.Context) (teams *[]models.Team, err error)
	GetTeamByName(ctx context.Context, name string) (team *models.Team, err error)
	GetTeamByNameLike(ctx context.Context, name string) (teams *[]models.Team, err error)
}

type TeamWriter interface {
	CreateTeam(ctx context.Context, team *models.Team) (err error)
	UpdateTeam(ctx context.Context, team *models.Team) (updatedTeam *models.Team, err error)
}

type TeamDeleter interface {
	DeleteTeam(ctx context.Context, teamId string) (err error)
}
