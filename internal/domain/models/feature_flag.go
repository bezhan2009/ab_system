package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeatureFlag struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Key          string `gorm:"uniqueIndex;not null"`
	Type         string `gorm:"column:type;size:20;not null;check:type IN ('string','number','boolean')"`
	DefaultValue string `gorm:"column:default_value;not null"`

	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	User        User      `gorm:"foreignKey:UserID"`
	Description string
}
