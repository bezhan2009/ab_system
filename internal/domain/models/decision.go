package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Decision struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       string         `gorm:"index;not null"`
	ExperimentID *uuid.UUID     `gorm:"type:uuid;index"`
	VariantID    *uuid.UUID     `gorm:"type:uuid;index"`
	Attributes   datatypes.JSON `json:"attributes"`
	FlagKey      string         `gorm:"not null"`
	Value        string         `gorm:"not null"`
	CreatedAt    time.Time      `gorm:"index"`
	ExpiresAt    time.Time      `gorm:"index"`
}
