package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationSettingsRepository struct {
	db *gorm.DB
}

func NewNotificationSettingsRepository(db *gorm.DB) *NotificationSettingsRepository {
	return &NotificationSettingsRepository{db: db}
}

func (r *NotificationSettingsRepository) GetNotificationSettingsByExperimentID(ctx context.Context, experimentID string) (settings *models.NotificationSettings, err error) {
	const op = "repository.postgres.GetNotificationSettingsByExperimentID"

	expID, err := uuid.Parse(experimentID)
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).
		Where("experiment_id = ?", expID).
		First(&settings).Error

	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return settings, nil
}

func (r *NotificationSettingsRepository) GetSlackWebhooksForExpNotification(ctx context.Context, experimentID string) (webhooks []string, err error) {
	const op = "repository.postgres.GetSlackWebhooksForExpNotification"

	expID, err := uuid.Parse(experimentID)
	if err != nil {
		return nil, TranslateGormError(err)
	}

	var webhooksJSON string
	err = r.db.WithContext(ctx).
		Raw("SELECT slack_webhooks::text FROM notification_settings WHERE experiment_id = ?", expID).
		Scan(&webhooksJSON).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	if webhooksJSON != "" {
		if err = json.Unmarshal([]byte(webhooksJSON), &webhooks); err != nil {
			logger.Error.Printf("[%s] TraceId=%s Error unmarshaling slack_webhooks: %v",
				op, observability.GetTraceID(ctx), err)

			return nil, TranslateGormError(err)
		}
	}

	return webhooks, nil
}

func (r *NotificationSettingsRepository) CreateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) (err error) {
	const op = "repository.postgres.CreateNotificationSettings"

	if err = r.db.WithContext(ctx).Create(settings).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *NotificationSettingsRepository) UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) (*models.NotificationSettings, error) {
	const op = "repository.postgres.UpdateNotificationSettings"

	if err := r.db.WithContext(ctx).
		Model(&models.NotificationSettings{}).
		Where("experiment_id = ?", settings.ExperimentID).
		Updates(map[string]interface{}{
			"chat_ids":   settings.ChatIds,
			"slack_webhooks": settings.SlackWebhooks,
			"updated_at": time.Now(),
		}).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)
		return nil, TranslateGormError(err)
	}

	var updated models.NotificationSettings
	if err := r.db.WithContext(ctx).
		Where("experiment_id = ?", settings.ExperimentID).
		First(&updated).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error reloading: %v",
			op, observability.GetTraceID(ctx), err)
		return nil, TranslateGormError(err)
	}

	return &updated, nil
}

func (r *NotificationSettingsRepository) DeleteNotificationSettings(ctx context.Context, settingsId string) (err error) {
	const op = "repository.postgres.DeleteNotificationSettings"

	uid, err := uuid.Parse(settingsId)
	if err != nil {
		return err
	}

	if err = r.db.WithContext(ctx).Where("id = ?", uid).Delete(&models.NotificationSettings{}).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}
