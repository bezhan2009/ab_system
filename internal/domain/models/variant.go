package models

import "github.com/google/uuid"

type Variant struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ExperimentID uuid.UUID `gorm:"not null;index"`
	Title        string    `gorm:"size:50;not null"`
	Value        string    `gorm:"not null"`
	Weight       int       `gorm:"not null"`
	IsControl    bool      `gorm:"not null;default:false"`
	Description  string    `gorm:"size:500"`
}
