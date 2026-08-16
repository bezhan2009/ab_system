package dto

import (
	"ab_system/internal/domain/models"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventType struct {
	ID                 string      `json:"id"`
	Title              string      `json:"title"`
	Description        string      `json:"description"`
	Schema             interface{} `json:"schema"`
	RequiresDecisionID bool        `json:"requires_decision_id"`
	RequiresUserID     bool        `json:"requires_user_id"`
	RequiresExposure   bool        `json:"requires_exposure"`
	CreatedAt          time.Time   `json:"createdAt"`
	UpdatedAt          time.Time   `json:"updatedAt"`
}

func (d *EventType) ToModel() (*models.EventType, error) {
	id := uuid.UUID{}
	if d.ID != "" {
		var err error
		id, err = uuid.Parse(d.ID)
		if err != nil {
			return nil, err
		}
	}

	schemaBytes, err := json.Marshal(d.Schema)
	if err != nil {
		return nil, err
	}

	return &models.EventType{
		ID:                 id,
		Title:              d.Title,
		Description:        d.Description,
		Schema:             schemaBytes,
		RequiresDecisionID: d.RequiresDecisionID,
		RequiresUserID:     d.RequiresUserID,
		RequiresExposure:   d.RequiresExposure,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}, nil
}

func (d *EventType) ToDTO(et *models.EventType) *EventType {
	var schema interface{}
	_ = json.Unmarshal(et.Schema, &schema)

	return &EventType{
		ID:                 et.ID.String(),
		Title:              et.Title,
		Description:        et.Description,
		Schema:             schema,
		RequiresDecisionID: et.RequiresDecisionID,
		RequiresUserID:     et.RequiresUserID,
		RequiresExposure:   et.RequiresExposure,
		CreatedAt:          et.CreatedAt,
		UpdatedAt:          et.UpdatedAt,
	}
}

func (d *EventType) ToDTOs(ets []models.EventType) []*EventType {
	result := make([]*EventType, len(ets))
	for i := range ets {
		result[i] = d.ToDTO(&ets[i])
	}
	return result
}
