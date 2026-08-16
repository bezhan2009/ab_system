package utils

import (
	notifydto "ab_system/pkg/notifications/dto"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateNotificationKey(notifyService string, req *notifydto.NotifyRequest) string {
	data := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s",
		notifyService,
		req.EventType,
		req.ExperimentID,
		req.Experiment,
		req.Status,
		req.Metric,
		req.Threshold,
	)

	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("notify:%s:%s", notifyService, hex.EncodeToString(hash[:]))
}

func GetEventTypeByExperimentStatus(status string) string {
	switch status {
	case "in_review":
		return "experiment_review"
	case "approved":
		return "experiment_approved"
	case "running":
		return "experiment_started"
	case "paused":
		return "experiment_paused"
	case "completed":
		return "experiment_completed"
	case "rejected":
		return "experiment_rejected"
	default:
		return ""
	}
}