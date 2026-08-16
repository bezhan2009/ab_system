package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type NotificationSettingsReader interface {
	GetNotificationSettingsByExperimentID(ctx context.Context, experimentID string) (settings *models.NotificationSettings, err error)
	GetSlackWebhooksForExpNotification(ctx context.Context, experimentID string) (webhooks []string, err error)
}

type NotificationSettingsWriter interface {
	CreateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) (err error)
	UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) (updated *models.NotificationSettings, err error)
}

type NotificationSettingsDeleter interface {
	DeleteNotificationSettings(ctx context.Context, settingsId string) (err error)
}
