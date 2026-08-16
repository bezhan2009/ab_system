package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type EventType struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"size:1000"`

	Schema datatypes.JSON `gorm:"type:json;not null"`

	RequiresDecisionID bool `gorm:"not null;default:true"`
	RequiresUserID     bool `gorm:"not null;default:true"`

	RequiresExposure bool `gorm:"not null;default:true"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
