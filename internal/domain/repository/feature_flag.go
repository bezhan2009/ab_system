package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type FeatureFlagReader interface {
	GetAllFeatureFlags(ctx context.Context) (featureFlags []models.FeatureFlag, err error)
	GetFeatureFlagById(ctx context.Context, featureFlagById string) (featureFlag models.FeatureFlag, err error)
	GetFeatureFlagByKey(ctx context.Context, key string) (featureFlag models.FeatureFlag, err error)
	GetFeatureFlagsByOwner(ctx context.Context, owner string) ([]models.FeatureFlag, error)
}

type FeatureFlagWriter interface {
	CreateFeatureFlag(ctx context.Context, featureFlag *models.FeatureFlag) (err error)
	UpdateFeatureFlag(ctx context.Context, featureFlag *models.FeatureFlag) (updatedFeatureFlag *models.FeatureFlag, err error)
}

type FeatureFlagDeleter interface {
	DeleteFeatureFlag(ctx context.Context, featureFlagById string) (err error)
}
