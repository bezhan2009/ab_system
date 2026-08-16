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

type MetricRepository struct {
	db *gorm.DB
}

type ExperimentMetricRepository struct {
	db *gorm.DB
}

func NewMetricRepository(db *gorm.DB) *MetricRepository {
	return &MetricRepository{db: db}
}

func NewExperimentMetricRepository(db *gorm.DB) *ExperimentMetricRepository {
	return &ExperimentMetricRepository{db: db}
}

func (r *MetricRepository) GetMetricByID(ctx context.Context, id string) (metric *models.Metric, err error) {
	const op = "repository.postgres.GetMetricByID"

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var m models.Metric
	if err = r.db.WithContext(ctx).Where("id = ?", uid).First(&m).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}
	return &m, nil
}

func (r *MetricRepository) GetAllMetrics(ctx context.Context) (metrics []models.Metric, err error) {
	const op = "repository.postgres.GetAllMetrics"

	if err = r.db.WithContext(ctx).Find(&metrics).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}
	return metrics, nil
}

func (r *MetricRepository) GetMetricByTitle(ctx context.Context, title string) (metric *models.Metric, err error) {
	const op = "repository.postgres.GetMetricByTitle"

	var m models.Metric
	if err = r.db.WithContext(ctx).Where("title = ?", title).First(&m).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return &m, nil
}

func (r *MetricRepository) CreateMetric(ctx context.Context, metric *models.Metric) (err error) {
	const op = "repository.postgres.CreateMetric"

	if err = r.db.WithContext(ctx).Create(metric).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *MetricRepository) UpdateMetric(ctx context.Context, metric *models.Metric) (updatedMetric *models.Metric, err error) {
	const op = "repository.postgres.UpdateMetric"

	if err = r.db.WithContext(ctx).Clauses(clause.Returning{}).Updates(metric).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return metric, nil
}

func (r *MetricRepository) DeleteMetric(ctx context.Context, id string) (err error) {
	const op = "repository.postgres.DeleteMetric"

	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err = r.db.WithContext(ctx).Where("id = ?", uid).Delete(&models.Metric{}).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *ExperimentMetricRepository) GetExperimentMetrics(ctx context.Context, experimentID string) (experimentMetrics []models.ExperimentMetric, err error) {
	const op = "repository.postgres.GetExperimentMetrics"

	expUID, err := uuid.Parse(experimentID)
	if err != nil {
		return nil, err
	}

	if err = r.db.WithContext(ctx).Where("experiment_id = ?", expUID).Find(&experimentMetrics).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experimentMetrics, nil
}

func (r *ExperimentMetricRepository) GetGuardrailsForExperiment(ctx context.Context, experimentID string) (guardrails []models.ExperimentMetric, err error) {
	const op = "repository.postgres.GetGuardrailsForExperiment"

	expUID, err := uuid.Parse(experimentID)
	if err != nil {
		return nil, err
	}

	if err = r.db.WithContext(ctx).
		Where("experiment_id = ? AND is_guardrail = ?", expUID, true).
		Find(&guardrails).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return guardrails, nil
}

func (r *ExperimentMetricRepository) AddMetricToExperiment(ctx context.Context, metricBinding *models.ExperimentMetric) (err error) {
	const op = "repository.postgres.AddMetricToExperiment"

	if err = r.db.WithContext(ctx).Create(metricBinding).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *ExperimentMetricRepository) UpdateExperimentMetric(ctx context.Context, metricBinding *models.ExperimentMetric) (updatedBinding *models.ExperimentMetric, err error) {
	const op = "repository.postgres.UpdateExperimentMetric"

	if err = r.db.WithContext(ctx).
		Where("experiment_id = ? AND metric_id = ?", metricBinding.ExperimentID, metricBinding.MetricID).
		Select("is_primary", "is_guardrail", "threshold", "operator", "window_min", "action").
		Clauses(clause.Returning{}).
		Updates(metricBinding).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return metricBinding, nil
}

func (r *ExperimentMetricRepository) RemoveMetricFromExperiment(ctx context.Context, experimentID, metricID string) (err error) {
	const op = "repository.postgres.RemoveMetricFromExperiment"

	if err = r.db.WithContext(ctx).
		Where("experiment_id = ? AND metric_id = ?", experimentID, metricID).
		Delete(&models.ExperimentMetric{}).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}
