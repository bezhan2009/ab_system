package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, teamId string) (team *models.Team, err error) {
	const op = "repository.postgres.GetTeamByID"

	uid, err := uuid.Parse(teamId)
	if err != nil {
		return nil, err
	}

	if err = r.db.WithContext(ctx).Where("id = ?", uid).First(&team).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return team, nil
}

func (r *TeamRepository) GetAllTeams(ctx context.Context) (teams *[]models.Team, err error) {
	const op = "repository.postgres.GetAllTeams"

	if err = r.db.WithContext(ctx).Find(&teams).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return teams, nil
}

func (r *TeamRepository) GetTeamByName(ctx context.Context, name string) (team *models.Team, err error) {
	const op = "repository.postgres.GetTeamByName"

	if err = r.db.WithContext(ctx).Where("name = ?", name).First(&team).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return team, nil
}

func (r *TeamRepository) GetTeamByNameLike(ctx context.Context, name string) (teams *[]models.Team, err error) {
	const op = "repository.postgres.GetTeamByNameLike"

	if err = r.db.WithContext(ctx).Where("name like ?", "%"+name+"%").Find(&teams).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return teams, nil
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *models.Team) (err error) {
	const op = "repository.postgres.CreateTeam"

	if err = r.db.WithContext(ctx).Create(team).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *TeamRepository) UpdateTeam(ctx context.Context, team *models.Team) (updatedTeam *models.Team, err error) {
	const op = "repository.postgres.UpdateTeam"

	if err = r.db.WithContext(ctx).Clauses(clause.Returning{}).Updates(team).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while updating team: %s",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return team, nil
}

func (r *TeamRepository) DeleteTeam(ctx context.Context, teamId string) (err error) {
	const op = "repository.postgres.DeleteTeam"

	uid, err := uuid.Parse(teamId)
	if err != nil {
		return err
	}

	if err = r.db.WithContext(ctx).Delete(&models.Team{}, "id = ?", uid).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}
