package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Event struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	EventID string `gorm:"size:255;not null;uniqueIndex"`

	EventTypeID string    `gorm:"size:100;not null;index"`
	EventType   EventType `gorm:"foreignKey:EventTypeID"`

	DecisionID string `gorm:"size:255;index;not null"`
	UserID     string `gorm:"size:255;index;not null"`

	Payload datatypes.JSON `gorm:"type:json"`

	ReceivedAt time.Time  `gorm:"not null;index"`
	ClientTime *time.Time `gorm:"index"`
}
