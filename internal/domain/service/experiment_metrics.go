package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/errs"
	"context"

	"github.com/google/uuid"
)

type ExperimentMetricService struct {
	expMetricReader  repository.ExperimentMetricReader
	expMetricWriter  repository.ExperimentMetricWriter
	expMetricDeleter repository.ExperimentMetricDeleter
	experimentReader repository.ExperimentReader
	metricReader     repository.MetricReader
}

func NewExperimentMetricService(
	expMetricReader repository.ExperimentMetricReader,
	expMetricWriter repository.ExperimentMetricWriter,
	expMetricDeleter repository.ExperimentMetricDeleter,
	experimentReader repository.ExperimentReader,
	metricReader repository.MetricReader,
) *ExperimentMetricService {
	return &ExperimentMetricService{
		expMetricReader:  expMetricReader,
		expMetricWriter:  expMetricWriter,
		expMetricDeleter: expMetricDeleter,
		experimentReader: experimentReader,
		metricReader:     metricReader,
	}
}

func (s *ExperimentMetricService) GetExperimentMetrics(ctx context.Context, experimentID string) ([]models.ExperimentMetric, error) {
	_, err := s.experimentReader.GetExperimentByID(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	metrics, err := s.expMetricReader.GetExperimentMetrics(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (s *ExperimentMetricService) GetGuardrailsForExperiment(ctx context.Context, experimentID string) ([]models.ExperimentMetric, error) {
	_, err := s.experimentReader.GetExperimentByID(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	guardrails, err := s.expMetricReader.GetGuardrailsForExperiment(ctx, experimentID)
	if err != nil {
		return nil, err
	}
	return guardrails, nil
}

func (s *ExperimentMetricService) AddMetricToExperiment(ctx context.Context, metricBinding *models.ExperimentMetric) error {
	_, err := s.experimentReader.GetExperimentByID(ctx, metricBinding.ExperimentID.String())
	if err != nil {
		return err
	}

	_, err = s.metricReader.GetMetricByID(ctx, metricBinding.MetricID.String())
	if err != nil {
		return err
	}

	existing, err := s.expMetricReader.GetExperimentMetrics(ctx, metricBinding.ExperimentID.String())
	if err != nil {
		return err
	}

	for _, em := range existing {
		if em.MetricID == metricBinding.MetricID {
			return errs.ErrAlreadyExists
		}
	}

	if metricBinding.IsPrimary {
		if err = s.ensureSinglePrimary(ctx, metricBinding.ExperimentID.String(), metricBinding.MetricID); err != nil {
			return err
		}
	}

	return s.expMetricWriter.AddMetricToExperiment(ctx, metricBinding)
}

func (s *ExperimentMetricService) UpdateExperimentMetric(ctx context.Context, metricBinding *models.ExperimentMetric) (*models.ExperimentMetric, error) {
	existing, err := s.expMetricReader.GetExperimentMetrics(ctx, metricBinding.ExperimentID.String())
	if err != nil {
		return nil, err
	}

	found := false
	for _, em := range existing {
		if em.MetricID == metricBinding.MetricID {
			found = true
			break
		}
	}

	if !found {
		return nil, errs.ErrRecordNotFound
	}

	if metricBinding.IsPrimary {
		if err = s.ensureSinglePrimary(ctx, metricBinding.ExperimentID.String(), metricBinding.MetricID); err != nil {
			return nil, err
		}
	}

	updated, err := s.expMetricWriter.UpdateExperimentMetric(ctx, metricBinding)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *ExperimentMetricService) RemoveMetricFromExperiment(ctx context.Context, experimentID, metricID string) error {
	return s.expMetricDeleter.RemoveMetricFromExperiment(ctx, experimentID, metricID)
}

func (s *ExperimentMetricService) ensureSinglePrimary(ctx context.Context, experimentID string, currentMetricID uuid.UUID) error {
	metrics, err := s.expMetricReader.GetExperimentMetrics(ctx, experimentID)
	if err != nil {
		return err
	}

	for _, em := range metrics {
		if em.IsPrimary && em.MetricID != currentMetricID {
			return errs.ErrMultiplePrimaryMetrics
		}
	}

	return nil
}
