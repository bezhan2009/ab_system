package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type NotificationSettings struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	ExperimentID *uuid.UUID  `gorm:"size:255;index"`
	Experiment   *Experiment `gorm:"foreignkey:ExperimentID"`

	ChatIds       datatypes.JSON `gorm:"type:json;default:'[]'"`
	SlackWebhooks datatypes.JSON `gorm:"type:json;default:'[]'"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
