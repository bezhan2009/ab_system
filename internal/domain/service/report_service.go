package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/internal/http/dto"
	"ab_system/internal/lib/metrics"
	"context"
	"time"
)

type ReportService struct {
	eventReader     repository.EventReader
	decisionReader  repository.DecisionReader
	metricReader    repository.MetricReader
	expMetricReader repository.ExperimentMetricReader
	variantReader   repository.VariantReader

	metricLib *metrics.MetricLib
}

func NewReportService(
	eventReader repository.EventReader,
	decisionReader repository.DecisionReader,
	metricReader repository.MetricReader,
	expMetricReader repository.ExperimentMetricReader,
	variantReader repository.VariantReader,
	metricLib *metrics.MetricLib) *ReportService {
	return &ReportService{
		eventReader:     eventReader,
		decisionReader:  decisionReader,
		metricReader:    metricReader,
		expMetricReader: expMetricReader,
		variantReader:   variantReader,
		metricLib:       metricLib,
	}
}

func (s *ReportService) GetExperimentReport(ctx context.Context, experimentID string, from, to time.Time) (*dto.ExperimentReport, error) {
	expMetrics, err := s.expMetricReader.GetExperimentMetrics(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	decisions, err := s.decisionReader.GetDecisionsByExperimentAndTime(ctx, experimentID, from, to)
	if err != nil {
		return nil, err
	}

	variantMap := make(map[string][]models.Decision)
	for _, d := range decisions {
		if d.VariantID != nil {
			variantMap[d.VariantID.String()] = append(variantMap[d.VariantID.String()], d)
		}
	}

	var variantReports []dto.VariantReport
	for vid, vDecisions := range variantMap {
		decIDs := make([]string, len(vDecisions))
		for i, d := range vDecisions {
			decIDs[i] = d.ID.String()
		}

		var metricValues []dto.MetricValue
		for _, em := range expMetrics {
			metric, err := s.metricReader.GetMetricByID(ctx, em.MetricID.String())
			if err != nil {
				continue
			}

			val, err := s.metricLib.CalculateMetric(ctx, metric, decIDs, from, to, true)
			if err != nil {
				continue
			}

			displayValue := val
			if metric.Type == models.MetricTypeRatio && metric.Unit == "%" {
				displayValue = val * 100
			}
			metricValues = append(metricValues, dto.MetricValue{
				MetricID:    metric.ID.String(),
				MetricTitle: metric.Title,
				Value:       displayValue,
				Unit:        metric.Unit,
			})
		}

		variantName := ""
		if len(vDecisions) > 0 {
			variant, err := s.variantReader.GetVariantById(ctx, vDecisions[0].VariantID.String())
			if err != nil {
				continue
			}

			variantName = variant.Title
		}

		variantReports = append(variantReports, dto.VariantReport{
			VariantID:    vid,
			VariantName:  variantName,
			MetricValues: metricValues,
		})
	}

	return &dto.ExperimentReport{
		ExperimentID: experimentID,
		From:         from,
		To:           to,
		Variants:     variantReports,
	}, nil
}
