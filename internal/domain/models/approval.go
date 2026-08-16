package models

import (
	"time"

	"github.com/google/uuid"
)

type Approval struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ExperimentID      uuid.UUID `gorm:"type:uuid;not null;index"`
	ExperimentVersion int       `gorm:"not null;default:0"`
	ApproverID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Approved          bool      `gorm:"not null"`
	Comment           string    `gorm:"size:1000"`
	CreatedAt         time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
