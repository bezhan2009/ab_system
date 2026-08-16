package dto

import (
	"ab_system/internal/domain/models"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NotificationSettings struct {
	ID           string `json:"id"`
	ExperimentID string `json:"experiment_id"`

	ChatIds       []int64  `json:"chat_ids"`
	SlackWebhooks []string `json:"slack_webhooks"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (d *NotificationSettings) ToModel() (*models.NotificationSettings, error) {
	id := uuid.UUID{}
	if d.ID != "" {
		var err error
		id, err = uuid.Parse(d.ID)
		if err != nil {
			return nil, err
		}
	}

	expID, err := uuid.Parse(d.ExperimentID)
	if err != nil {
		return nil, err
	}

	chatIdsJSON, err := json.Marshal(d.ChatIds)
	if err != nil {
		return nil, err
	}

	slackWebhooksJSON, err := json.Marshal(d.SlackWebhooks)
	if err != nil {
		return nil, err
	}

	return &models.NotificationSettings{
		ID:            id,
		ExperimentID:  &expID,
		ChatIds:       chatIdsJSON,
		SlackWebhooks: slackWebhooksJSON,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}, nil
}

func (d *NotificationSettings) ToDTO(model *models.NotificationSettings) (*NotificationSettings, error) {
	var chatIds []int64
	if err := json.Unmarshal(model.ChatIds, &chatIds); err != nil {
		return nil, err
	}

	var slackWebhooks []string
	if err := json.Unmarshal(model.SlackWebhooks, &slackWebhooks); err != nil {
		return nil, err
	}

	return &NotificationSettings{
		ID:            model.ID.String(),
		ExperimentID:  model.ExperimentID.String(),
		ChatIds:       chatIds,
		SlackWebhooks: slackWebhooks,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (d *NotificationSettings) ToDTOs(models []models.NotificationSettings) ([]*NotificationSettings, error) {
	result := make([]*NotificationSettings, len(models))
	for i := range models {
		dto, err := d.ToDTO(&models[i])
		if err != nil {
			return nil, err
		}
		result[i] = dto
	}

	return result, nil
}
