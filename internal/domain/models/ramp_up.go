package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ExperimentRampUp struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	ExperimentID *uuid.UUID
	Experiment   *Experiment `gorm:"foreignkey:ExperimentID"`

	RampEnabled         bool           `gorm:"default:false"`
	RampSteps           datatypes.JSON `gorm:"type:json;default:'[]'"`
	RampCurrentStep     int            `gorm:"default:0"`
	RampLastIncrease    time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	RampIntervalMinutes int            `gorm:"default:60"`
}

//{
//"ramp_enabled": true,
//"ramp_steps": [1, 5, 10, 25, 50, 100],
//"ramp_interval_minutes": 60
//}
