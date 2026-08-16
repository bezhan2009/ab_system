package models

import (
	"time"

	"github.com/google/uuid"
)

type ApproverGroup struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ExperimenterID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	MinApprovals   int       `gorm:"not null;default:1"`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type ApproverGroupMember struct {
	ApproverGroupID uuid.UUID `gorm:"type:uuid;primaryKey"`
	ApproverID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt       time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
