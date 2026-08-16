package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"

	"gorm.io/gorm"
)

type GuardrailTriggerRepository struct {
	db *gorm.DB
}

func NewGuardrailTriggerRepository(db *gorm.DB) *GuardrailTriggerRepository {
	return &GuardrailTriggerRepository{
		db: db,
	}
}

func (r *GuardrailTriggerRepository) CreateTrigger(ctx context.Context, trigger *models.GuardrailTrigger) (err error) {
	const op = "repository.postgres.GuardrailTriggerRepository.CreateTrigger"

	if err = r.db.WithContext(ctx).Create(trigger).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}
	return nil
}
