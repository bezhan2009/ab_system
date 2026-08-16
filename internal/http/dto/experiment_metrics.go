package dto

import (
	"ab_system/internal/domain/models"

	"github.com/google/uuid"
)

type ExperimentMetric struct {
	ExperimentID string   `json:"experiment_id"`
	MetricID     string   `json:"metric_id"`
	IsPrimary    bool     `json:"is_primary"`
	IsGuardrail  bool     `json:"is_guardrail"`
	Threshold    *float64 `json:"threshold,omitempty"`
	Operator     *string  `json:"operator,omitempty" binding:"omitempty,oneof='>' '>=' '<' '<='"`
	WindowMin    *int     `json:"window_min,omitempty"`
	Action       *string  `json:"action,omitempty" binding:"omitempty,oneof=pause rollback"`
}

func (d *ExperimentMetric) ToModel() (*models.ExperimentMetric, error) {
	expID, err := uuid.Parse(d.ExperimentID)
	if err != nil {
		return nil, err
	}
	metID, err := uuid.Parse(d.MetricID)
	if err != nil {
		return nil, err
	}

	model := &models.ExperimentMetric{
		ExperimentID: expID,
		MetricID:     metID,
		IsPrimary:    d.IsPrimary,
		IsGuardrail:  d.IsGuardrail,
	}

	if d.Threshold != nil {
		model.Threshold = *d.Threshold
	}
	if d.Operator != nil {
		model.Operator = *d.Operator
	}
	if d.WindowMin != nil {
		model.WindowMin = *d.WindowMin
	}
	if d.Action != nil {
		model.Action = *d.Action
	}

	return model, nil
}

func (d *ExperimentMetric) ToDTO(model *models.ExperimentMetric) *ExperimentMetric {
	dto := &ExperimentMetric{
		ExperimentID: model.ExperimentID.String(),
		MetricID:     model.MetricID.String(),
		IsPrimary:    model.IsPrimary,
		IsGuardrail:  model.IsGuardrail,
		Threshold:    &model.Threshold,
		Operator:     &model.Operator,
		WindowMin:    &model.WindowMin,
		Action:       &model.Action,
	}
	return dto
}

func (d *ExperimentMetric) ToDTOs(models []models.ExperimentMetric) []*ExperimentMetric {
	result := make([]*ExperimentMetric, len(models))
	for i := range models {
		result[i] = d.ToDTO(&models[i])
	}
	return result
}
