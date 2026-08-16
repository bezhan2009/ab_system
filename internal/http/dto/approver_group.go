package dto

import (
	"ab_system/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type ApproverGroup struct {
	ID             string    `json:"id"`
	ExperimenterID string    `json:"experimenter_id" binding:"required"`
	MinApprovals   int       `json:"min_approvals" binding:"required,min=1"`
	ApproverIDs    []string  `json:"approver_ids" binding:"required,min=1"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (d *ApproverGroup) ToModel() (*models.ApproverGroup, []uuid.UUID, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil && d.ID != "" {
		return nil, nil, err
	}

	expID, err := uuid.Parse(d.ExperimenterID)
	if err != nil {
		return nil, nil, err
	}

	approverUUIDs := make([]uuid.UUID, 0, len(d.ApproverIDs))
	for _, aid := range d.ApproverIDs {
		uid, err := uuid.Parse(aid)
		if err != nil {
			return nil, nil, err
		}

		approverUUIDs = append(approverUUIDs, uid)
	}

	return &models.ApproverGroup{
		ID:             id,
		ExperimenterID: expID,
		MinApprovals:   d.MinApprovals,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}, approverUUIDs, nil
}

func (d *ApproverGroup) ToDTO(group *models.ApproverGroup, approvers []models.User) *ApproverGroup {
	approverIDs := make([]string, len(approvers))
	for i, u := range approvers {
		approverIDs[i] = u.ID.String()
	}

	return &ApproverGroup{
		ID:             group.ID.String(),
		ExperimenterID: group.ExperimenterID.String(),
		MinApprovals:   group.MinApprovals,
		ApproverIDs:    approverIDs,
		CreatedAt:      group.CreatedAt,
		UpdatedAt:      group.UpdatedAt,
	}
}
