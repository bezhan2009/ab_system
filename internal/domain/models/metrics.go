package models

import (
	"time"

	"github.com/google/uuid"
)

type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"   // просто счётчик событий
	MetricTypeRatio     MetricType = "ratio"     // доля (числитель/знаменатель)
	MetricTypeHistogram MetricType = "histogram" // для latency (среднее, p95)
)

type Metric struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string     `gorm:"size:255;not null;uniqueIndex"`
	Description string     `gorm:"size:1000"`
	Type        MetricType `gorm:"size:50;not null"`

	CounterEventType string `gorm:"size:255"`

	// Для ratio
	NumeratorEventType   string `gorm:"size:255"`
	DenominatorEventType string `gorm:"size:255"`

	// Для histogram
	HistogramEventType string `gorm:"size:255"`
	HistogramField     string `gorm:"size:255"`

	RequiresExposure bool `gorm:"not null;default:true"`

	Unit string `gorm:"size:50"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type ExperimentMetric struct {
	ExperimentID uuid.UUID `gorm:"type:uuid;primaryKey"`
	MetricID     uuid.UUID `gorm:"type:uuid;primaryKey"`

	IsPrimary   bool `gorm:"not null;default:false"`
	IsGuardrail bool `gorm:"not null;default:false"`

	// Для guardrail
	Threshold float64 `gorm:""`
	Operator  string  `gorm:"size:10"`
	WindowMin int     `gorm:""`
	Action    string  `gorm:"size:50"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type GuardrailTrigger struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ExperimentID uuid.UUID `gorm:"type:uuid;not null;index"`
	MetricID     uuid.UUID `gorm:"type:uuid;not null"`
	Threshold    float64
	Operator     string
	WindowMin    int
	ActualValue  float64
	Action       string
	TriggeredAt  time.Time
}
