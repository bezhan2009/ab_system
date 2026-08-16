package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/errs"
	"context"
	"errors"
)

type MetricService struct {
	metricReader  repository.MetricReader
	metricWriter  repository.MetricWriter
	metricDeleter repository.MetricDeleter
}

func NewMetricService(
	metricReader repository.MetricReader,
	metricWriter repository.MetricWriter,
	metricDeleter repository.MetricDeleter,
) *MetricService {
	return &MetricService{
		metricReader:  metricReader,
		metricWriter:  metricWriter,
		metricDeleter: metricDeleter,
	}
}

func (s *MetricService) GetAllMetrics(ctx context.Context) (metrics []models.Metric, err error) {
	metrics, err = s.metricReader.GetAllMetrics(ctx)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (s *MetricService) GetMetricByID(ctx context.Context, id string) (metric *models.Metric, err error) {
	metric, err = s.metricReader.GetMetricByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return metric, nil
}

func (s *MetricService) GetMetricByTitle(ctx context.Context, title string) (metric *models.Metric, err error) {
	metric, err = s.metricReader.GetMetricByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	return metric, nil
}

func (s *MetricService) CreateMetric(ctx context.Context, metric *models.Metric) (err error) {
	existing, err := s.metricReader.GetMetricByTitle(ctx, metric.Title)
	if err != nil && !errors.Is(err, errs.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return errs.ErrAlreadyExists
	}

	if err = s.validateMetric(metric); err != nil {
		return err
	}

	err = s.metricWriter.CreateMetric(ctx, metric)
	if err != nil {
		return err
	}
	return nil
}

func (s *MetricService) UpdateMetric(ctx context.Context, metric *models.Metric) (updatedMetric *models.Metric, err error) {
	existing, err := s.metricReader.GetMetricByID(ctx, metric.ID.String())
	if err != nil {
		return nil, err
	}

	if metric.Title != "" && metric.Title != existing.Title {
		dup, err := s.metricReader.GetMetricByTitle(ctx, metric.Title)
		if err != nil && !errors.Is(err, errs.ErrRecordNotFound) {
			return nil, err
		}
		if dup != nil && dup.ID != metric.ID {
			return nil, errs.ErrAlreadyExists
		}
	}

	if err = s.validateMetric(metric); err != nil {
		return nil, err
	}

	updatedMetric, err = s.metricWriter.UpdateMetric(ctx, metric)
	if err != nil {
		return nil, err
	}

	return updatedMetric, nil
}

func (s *MetricService) DeleteMetric(ctx context.Context, id string) (err error) {
	err = s.metricDeleter.DeleteMetric(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *MetricService) validateMetric(metric *models.Metric) error {
	switch metric.Type {
	case models.MetricTypeCounter:
		if metric.CounterEventType == "" {
			return errs.ErrInvalidMetricConfig
		}
	case models.MetricTypeRatio:
		if metric.NumeratorEventType == "" || metric.DenominatorEventType == "" {
			return errs.ErrInvalidMetricConfig
		}
	case models.MetricTypeHistogram:
		if metric.HistogramEventType == "" || metric.HistogramField == "" {
			return errs.ErrInvalidMetricConfig
		}
	default:
		return errs.ErrInvalidMetricType
	}

	return nil
}
