package postgres

import (
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationSettingsRepository struct {
	db *gorm.DB
}

func NewNotificationSettings(db *gorm.DB) *NotificationSettingsRepository {
	return &NotificationSettingsRepository{
		db: db,
	}
}

func (r *NotificationSettingsRepository) GetChatsForExpNotification(ctx context.Context, experimentID string) ([]int64, error) {
	const op = "repository.postgres.GetChatsForExpNotification"

	expID, err := uuid.Parse(experimentID)
	if err != nil {
		return nil, err
	}

	var chatIdsJSON string
	err = r.db.WithContext(ctx).
		Raw("SELECT chat_ids::text FROM notification_settings WHERE experiment_id = ?", expID).
		Scan(&chatIdsJSON).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error.Printf("[%s] TraceId=%s Error: %v", op, observability.GetTraceID(ctx), err)
		return nil, err
	}

	var chatIds []int64
	if chatIdsJSON != "" {
		if err = json.Unmarshal([]byte(chatIdsJSON), &chatIds); err != nil {
			logger.Error.Printf("[%s] TraceId=%s Error unmarshaling chat_ids: %v", op, observability.GetTraceID(ctx), err)
			return nil, err
		}
	}

	return chatIds, nil
}
