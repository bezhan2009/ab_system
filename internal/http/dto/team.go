package dto

import (
	"ab_system/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TeamMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	TeamID string `json:"team_id" binding:"required"`
}

func (d *Team) ToModel() (*models.Team, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil && d.ID != "" {
		return nil, err
	}
	return &models.Team{
		ID:          id,
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}, nil
}

func (d *Team) ToDTO(team *models.Team) *Team {
	return &Team{
		ID:          team.ID.String(),
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}
}

func (d *Team) ToDTOs(teams []models.Team) []*Team {
	result := make([]*Team, len(teams))
	for i := range teams {
		result[i] = d.ToDTO(&teams[i])
	}
	return result
}
