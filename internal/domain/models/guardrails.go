package models

import "github.com/google/uuid"

type Guardrails struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	EventTypeID uuid.UUID `gorm:"size:100;not null;index"`
	EventType   EventType `gorm:"foreignKey:EventTypeID"`

	ExperimentID uuid.UUID  `gorm:"type:uuid;not null;index"`
	Experiment   Experiment `gorm:"foreignKey:ExperimentID"`
}
