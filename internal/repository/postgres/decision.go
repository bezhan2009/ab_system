package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DecisionRepository struct {
	db *gorm.DB
}

func NewDecisionRepository(db *gorm.DB) *DecisionRepository {
	return &DecisionRepository{
		db: db,
	}
}

func (r *DecisionRepository) GetAllExperimentDecisions(ctx context.Context, experimentId string) (decisions *models.Decision, err error) {
	const op = "repository.postgres.GetAllExperimentDecisions"

	if err = r.db.WithContext(ctx).Find(&decisions, "experiment_id = ?", experimentId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting all decisions by experiment id: %v",
			op, observability.GetTraceID(ctx), err)

		return decisions, TranslateGormError(err)
	}

	return decisions, nil
}

func (r *DecisionRepository) GetDecisionByUserId(ctx context.Context, userId string) (decision *models.Decision, err error) {
	const op = "repository.postgres.GetDecisionByUserId"

	if err = r.db.WithContext(ctx).First(&decision, "user_id = ?", userId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting decision by user id: %v",
			op, observability.GetTraceID(ctx), err)

		return decision, TranslateGormError(err)
	}

	return decision, nil
}

func (r *DecisionRepository) GetDecisionsByExperimentAndTime(ctx context.Context, experimentID string, from, to time.Time) (decisions []models.Decision, err error) {
	const op = "repository.postgres.DecisionRepository.GetDecisionsByExperimentAndTime"

	uid, err := uuid.Parse(experimentID)
	if err != nil {
		return nil, err
	}

	query := r.db.WithContext(ctx).
		Where("experiment_id = ?", uid).
		Where("created_at >= ? AND created_at < ?", from, to)

	if err = query.Find(&decisions).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v", op, observability.GetTraceID(ctx), err)
		return nil, TranslateGormError(err)
	}

	return decisions, nil
}

func (r *DecisionRepository) CreateDecision(ctx context.Context, decision *models.Decision) (err error) {
	const op = "repository.postgres.CreateDecision"

	if err = r.db.WithContext(ctx).Create(&decision).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while creating decision: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *DecisionRepository) UpdateDecision(ctx context.Context, decision *models.Decision) (updatedDecision *models.Decision, err error) {
	const op = "repository.postgres.UpdateDecision"

	if err = r.db.WithContext(ctx).Clauses(clause.Returning{}).Updates(decision).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while updating decision: %s",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return decision, nil
}

func (r *DecisionRepository) DeleteDecision(ctx context.Context, decisionId string) (err error) {
	const op = "repository.postgres.DeleteDecision"

	if err = r.db.WithContext(ctx).Delete(&models.Decision{}, "id = ?", decisionId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while deleting decision: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}
