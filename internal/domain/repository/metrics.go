package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type MetricReader interface {
	GetMetricByID(ctx context.Context, id string) (metric *models.Metric, err error)
	GetAllMetrics(ctx context.Context) (metrics []models.Metric, err error)
	GetMetricByTitle(ctx context.Context, title string) (metric *models.Metric, err error)
}

type MetricWriter interface {
	CreateMetric(ctx context.Context, metric *models.Metric) (err error)
	UpdateMetric(ctx context.Context, metric *models.Metric) (updatedMetric *models.Metric, err error)
}

type MetricDeleter interface {
	DeleteMetric(ctx context.Context, id string) (err error)
}

type GuardrailTriggerWriter interface {
	CreateTrigger(ctx context.Context, trigger *models.GuardrailTrigger) (err error)
}
