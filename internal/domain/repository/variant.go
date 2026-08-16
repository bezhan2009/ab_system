package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type VariantReader interface {
	GetAllExperimentVariants(ctx context.Context, experimentId string) (variants *[]models.Variant, err error)
	GetControlExperimentVariant(ctx context.Context, experimentId string) (variant *models.Variant, err error)
	GetVariantById(ctx context.Context, variantId string) (variant *models.Variant, err error)
}

type VariantWriter interface {
	CreateVariant(ctx context.Context, variant *models.Variant) (err error)
	UpdateVariant(ctx context.Context, variant *models.Variant) (updatedVariant *models.Variant, err error)
}

type VariantDeleter interface {
	DeleteVariant(ctx context.Context, variantId string) (err error)
}
