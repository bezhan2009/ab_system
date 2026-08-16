package dto

import (
	"ab_system/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type Approval struct {
	ID                string    `json:"id"`
	ExperimentID      string    `json:"experiment_id"`
	ExperimentVersion int       `json:"experiment_version"`
	ApproverID        string    `json:"approver_id"`
	Approved          bool      `json:"approved"`
	Comment           string    `json:"comment,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
}

type ApprovalRequest struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment,omitempty"`
}

func (d *Approval) ToModel() (*models.Approval, error) {
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

	apprID, err := uuid.Parse(d.ApproverID)
	if err != nil {
		return nil, err
	}

	return &models.Approval{
		ID:                id,
		ExperimentID:      expID,
		ExperimentVersion: d.ExperimentVersion,
		ApproverID:        apprID,
		Approved:          d.Approved,
		Comment:           d.Comment,
		CreatedAt:         d.CreatedAt,
	}, nil
}

func (d *Approval) ToDTO(approval *models.Approval) *Approval {
	return &Approval{
		ID:                approval.ID.String(),
		ExperimentID:      approval.ExperimentID.String(),
		ExperimentVersion: approval.ExperimentVersion,
		ApproverID:        approval.ApproverID.String(),
		Approved:          approval.Approved,
		Comment:           approval.Comment,
		CreatedAt:         approval.CreatedAt,
	}
}

func (d *Approval) ToDTOs(approvals []models.Approval) []*Approval {
	result := make([]*Approval, len(approvals))
	for i := range approvals {
		result[i] = d.ToDTO(&approvals[i])
	}
	return result
}
