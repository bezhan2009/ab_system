package dto

import (
	"ab_system/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`

	CounterEventType     string `json:"counter_event_type,omitempty"`
	NumeratorEventType   string `json:"numerator_event_type,omitempty"`
	DenominatorEventType string `json:"denominator_event_type,omitempty"`
	HistogramEventType   string `json:"histogram_event_type,omitempty"`
	HistogramField       string `json:"histogram_field,omitempty"`

	RequiresExposure bool   `json:"requires_exposure"`
	Unit             string `json:"unit"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (d *Metric) ToModel() (*models.Metric, error) {
	id := uuid.UUID{}
	if d.ID != "" {
		var err error
		id, err = uuid.Parse(d.ID)
		if err != nil {
			return nil, err
		}
	}
	return &models.Metric{
		ID:                   id,
		Title:                d.Title,
		Description:          d.Description,
		Type:                 models.MetricType(d.Type),
		CounterEventType:     d.CounterEventType,
		NumeratorEventType:   d.NumeratorEventType,
		DenominatorEventType: d.DenominatorEventType,
		HistogramEventType:   d.HistogramEventType,
		HistogramField:       d.HistogramField,
		RequiresExposure:     d.RequiresExposure,
		Unit:                 d.Unit,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}, nil
}

func (d *Metric) ToDTO(m *models.Metric) *Metric {
	return &Metric{
		ID:                   m.ID.String(),
		Title:                 m.Title,
		Description:          m.Description,
		Type:                 string(m.Type),
		CounterEventType:     m.CounterEventType,
		NumeratorEventType:   m.NumeratorEventType,
		DenominatorEventType: m.DenominatorEventType,
		HistogramEventType:   m.HistogramEventType,
		HistogramField:       m.HistogramField,
		RequiresExposure:     m.RequiresExposure,
		Unit:                 m.Unit,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func (d *Metric) ToDTOs(ms []models.Metric) []*Metric {
	res := make([]*Metric, len(ms))
	for i := range ms {
		res[i] = d.ToDTO(&ms[i])
	}
	return res
}
