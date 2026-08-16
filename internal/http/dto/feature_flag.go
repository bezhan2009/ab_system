package dto

import (
	"ab_system/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type FeatureFlag struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Key          string `json:"key"`
	Type         string `json:"type"`
	DefaultValue string `json:"default_value"`
	Description  string `json:"description"`

	UserID string `json:"user_id"`
}

func (d *FeatureFlag) ToModel() (*models.FeatureFlag, error) {
	userID, err := uuid.Parse(d.UserID)
	if err != nil {
		return nil, err
	}

	flagID, err := uuid.Parse(d.ID)
	if err != nil {
		flagID = uuid.UUID{}
	}

	return &models.FeatureFlag{
		ID:           flagID,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		Key:          d.Key,
		Type:         d.Type,
		DefaultValue: d.DefaultValue,
		Description:  d.Description,
		UserID:       userID,
	}, nil
}

func (d *FeatureFlag) ToDTO(featureFlagModel *models.FeatureFlag) *FeatureFlag {
	return &FeatureFlag{
		ID:           featureFlagModel.ID.String(),
		CreatedAt:    featureFlagModel.CreatedAt,
		UpdatedAt:    featureFlagModel.UpdatedAt,
		Key:          featureFlagModel.Key,
		Type:         featureFlagModel.Type,
		DefaultValue: featureFlagModel.DefaultValue,
		Description:  featureFlagModel.Description,
		UserID:       featureFlagModel.UserID.String(),
	}
}

func (d *FeatureFlag) ToDTOs(featureFlagModels []models.FeatureFlag) []*FeatureFlag {
	result := make([]*FeatureFlag, len(featureFlagModels))
	for i := range featureFlagModels {
		result[i] = d.ToDTO(&featureFlagModels[i])
	}

	return result
}
