package service

import (
	"ab_system/internal/domain/models"
	telegramv1 "ab_system/internal/telegram_bot/gen/telegram_notifications/v1"
	"ab_system/pkg/utils"
	"context"
	"time"
)

func SendTelegramNotification(ctx context.Context, Notify func(ctx context.Context, req *telegramv1.NotifyRequest) (err error), req *telegramv1.NotifyRequest) (err error) {
	err = Notify(ctx, req)

	return err
}

func SendTelegramNotificationByExperimentStatus(
	ctx context.Context,
	notifyFunc func(ctx context.Context, req *telegramv1.NotifyRequest) (err error),
	experiment *models.Experiment,
	ownerEmail string,
) (err error) {

	eventType := utils.GetEventTypeByExperimentStatus(string(experiment.Status))

	req := &telegramv1.NotifyRequest{
		EventType:    eventType,
		ExperimentId: experiment.ID.String(),
		Experiment:   experiment.Title,
		FlagKey:      experiment.FlagKey,
		UserId:       ownerEmail,
		Timestamp:    time.Now().Format(time.RFC3339),
	}

	if experiment.Status == models.StatusCompleted && experiment.Conclusion != "" {
		req.Status = experiment.Conclusion
	}

	return notifyFunc(ctx, req)
}
