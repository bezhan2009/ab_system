package metrics

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/errs"
	"context"
	"time"
)

type MetricLib struct {
	eventReader repository.EventReader
}

func NewMetricLib(eventReader repository.EventReader) *MetricLib {
	return &MetricLib{eventReader: eventReader}
}

func (s *MetricLib) CalculateMetric(
	ctx context.Context,
	metric *models.Metric,
	decisionIDs []string,
	from, to time.Time,
	useClientTime bool,
) (float64, error) {
	switch metric.Type {
	case models.MetricTypeCounter:
		count, err := s.eventReader.CountEventsByTypeAndDecisions(ctx, metric.CounterEventType, decisionIDs, from, to, useClientTime)
		if err != nil {
			return 0, err
		}
		return float64(count), nil
	case models.MetricTypeRatio:
		num, err := s.eventReader.CountEventsByTypeAndDecisions(ctx, metric.NumeratorEventType, decisionIDs, from, to, useClientTime)
		if err != nil {
			return 0, err
		}
		den, err := s.eventReader.CountEventsByTypeAndDecisions(ctx, metric.DenominatorEventType, decisionIDs, from, to, useClientTime)
		if err != nil {
			return 0, err
		}
		if den == 0 {
			return 0, nil
		}
		return float64(num) / float64(den), nil
	case models.MetricTypeHistogram:
		avg, err := s.eventReader.AvgFieldByTypeAndDecisions(ctx, metric.HistogramEventType, metric.HistogramField, decisionIDs, from, to, useClientTime)
		if err != nil {
			return 0, err
		}
		return avg, nil
	default:
		return 0, errs.NewValidationError("metric.type", "Неизвестный тип метрики: " + string(metric.Type))
	}
}
