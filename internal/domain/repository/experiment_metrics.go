package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type ExperimentMetricReader interface {
	GetExperimentMetrics(ctx context.Context, experimentID string) (experimentMetrics []models.ExperimentMetric, err error)
	GetGuardrailsForExperiment(ctx context.Context, experimentID string) (guardrails []models.ExperimentMetric, err error)
}

type ExperimentMetricWriter interface {
	AddMetricToExperiment(ctx context.Context, metricBinding *models.ExperimentMetric) (err error)
	UpdateExperimentMetric(ctx context.Context, metricBinding *models.ExperimentMetric) (updatedBinding *models.ExperimentMetric, err error)
}

type ExperimentMetricDeleter interface {
	RemoveMetricFromExperiment(ctx context.Context, experimentID, metricID string) (err error)
}
