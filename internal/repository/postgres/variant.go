package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VariantRepository struct {
	db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) *VariantRepository {
	return &VariantRepository{
		db: db,
	}
}

func (r *VariantRepository) GetAllExperimentVariants(ctx context.Context, experimentId string) (variants *[]models.Variant, err error) {
	const op = "repository.postgres.GetAllExperimentVariants"

	if err = r.db.WithContext(ctx).Find(&variants, "experiment_id = ?", experimentId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting all experiment variants %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return variants, nil
}

func (r *VariantRepository) GetControlExperimentVariant(ctx context.Context, experimentId string) (variant *models.Variant, err error) {
	const op = "repository.postgres.GetControlExperimentVariant"

	if err = r.db.WithContext(ctx).First(&variant, "experiment_id = ? AND is_control = true", experimentId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting control experiment variant %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return variant, nil
}

func (r *VariantRepository) GetVariantById(ctx context.Context, variantId string) (variant *models.Variant, err error) {
	const op = "repository.postgres.GetVariantById"

	if err = r.db.WithContext(ctx).First(&variant, "id = ?", variantId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting variant %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return variant, nil
}

func (r *VariantRepository) CreateVariant(ctx context.Context, variant *models.Variant) (err error) {
	const op = "repository.postgres.CreateVariant"

	if err = r.db.WithContext(ctx).Create(variant).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while creating variant %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *VariantRepository) UpdateVariant(ctx context.Context, variant *models.Variant) (updatedVariant *models.Variant, err error) {
	const op = "repository.postgres.UpdateVariant"

	if err = r.db.WithContext(ctx).Clauses(clause.Returning{}).Updates(variant).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while updating variant: %s",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return variant, nil
}

func (r *VariantRepository) DeleteVariant(ctx context.Context, variantId string) (err error) {
	const op = "repository.postgres.DeleteVariant"

	if err = r.db.WithContext(ctx).Delete(&models.Variant{}, "id = ?", variantId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while deleting variant: %s",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}
