package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FeatureFlagRepository struct {
	db *gorm.DB
}

func NewFeatureFlagRepository(db *gorm.DB) *FeatureFlagRepository {
	return &FeatureFlagRepository{
		db: db,
	}
}

func (r *FeatureFlagRepository) GetAllFeatureFlags(ctx context.Context) (featureFlags []models.FeatureFlag, err error) {
	const op = "repository.postgres.GetAllFeatureFlags"

	if err = r.db.WithContext(ctx).Find(&featureFlags).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting all feature flags: %s",
			op, observability.GetTraceID(ctx), err)

		return featureFlags, TranslateGormError(err)
	}

	return featureFlags, nil
}

func (r *FeatureFlagRepository) GetFeatureFlagById(ctx context.Context, featureFlagById string) (featureFlag models.FeatureFlag, err error) {
	const op = "repository.postgres.GetFeatureFlagById"

	if err = r.db.WithContext(ctx).Where("id = ?", featureFlagById).First(&featureFlag).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting feature flag by id: %s",
			op, observability.GetTraceID(ctx), err)

		return featureFlag, TranslateGormError(err)
	}

	return featureFlag, nil
}

func (r *FeatureFlagRepository) GetFeatureFlagByKey(ctx context.Context, key string) (featureFlag models.FeatureFlag, err error) {
	const op = "repository.postgres.GetFeatureFlagByKey"

	if err = r.db.WithContext(ctx).Where("key = ?", key).First(&featureFlag).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting feature flag by key: %s",
			op, observability.GetTraceID(ctx), err)

		return featureFlag, TranslateGormError(err)
	}

	return featureFlag, nil
}

func (r *FeatureFlagRepository) GetFeatureFlagsByOwner(ctx context.Context, owner string) (featureFlags []models.FeatureFlag, err error) {
	const op = "repository.postgres.GetFeatureFlagsByOwner"

	if err = r.db.WithContext(ctx).Where("user_id = ?", owner).Find(&featureFlags).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting feature flags by owner: %s",
			op, observability.GetTraceID(ctx), err)

		return featureFlags, TranslateGormError(err)
	}

	return featureFlags, nil
}

func (r *FeatureFlagRepository) CreateFeatureFlag(ctx context.Context, featureFlag *models.FeatureFlag) (err error) {
	const op = "repository.postgres.CreateFeatureFlag"

	if err = r.db.WithContext(ctx).Create(featureFlag).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while creating feature flag: %s",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *FeatureFlagRepository) UpdateFeatureFlag(ctx context.Context, featureFlag *models.FeatureFlag) (updatedFeatureFlag *models.FeatureFlag, err error) {
	const op = "repository.postgres.UpdateFeatureFlag"

	if err = r.db.WithContext(ctx).Clauses(clause.Returning{}).Updates(featureFlag).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while updating feature flag: %s",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return featureFlag, nil
}

func (r *FeatureFlagRepository) DeleteFeatureFlag(ctx context.Context, featureFlagId string) (err error) {
	const op = "repository.postgres.DeleteFeatureFlag"

	if err = r.db.WithContext(ctx).Where("id = ?", featureFlagId).Delete(&models.FeatureFlag{}).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while deleting feature flag by id: %s",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}
