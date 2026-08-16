package service

import (
	"ab_system/internal/domain/redis"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/internal/telegram_bot/repository/postgres"
	"ab_system/pkg/logger"
	notifydto "ab_system/pkg/notifications/dto"
	"ab_system/pkg/notifications/message"
	"ab_system/pkg/utils"
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	bot                  *tgbotapi.BotAPI
	notificationSettings *postgres.NotificationSettingsRepository
	redisCache           redis.RedisCacheRepository
	ttl                  time.Duration
}

func NewTelegramNotifier(bot *tgbotapi.BotAPI, notificationSettings *postgres.NotificationSettingsRepository, redisCache redis.RedisCacheRepository, ttl time.Duration) *TelegramNotifier {
	return &TelegramNotifier{
		bot:                  bot,
		notificationSettings: notificationSettings,
		redisCache:           redisCache,
		ttl:                  ttl,
	}
}

func (s *TelegramNotifier) Send(ctx context.Context, req notifydto.NotifyRequest) error {
	const op = "service.TelegramNotifier.Send"

	key := utils.GenerateNotificationKey("telegram", &req)

	success, err := s.redisCache.SetNX(ctx, key, "1", s.ttl)
	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Cache error: %v", op, observability.GetTraceID(ctx), err)
	} else if !success {
		logger.Info.Printf("[%s] TraceId=%s Duplicate notification skipped, key: %s", op, observability.GetTraceID(ctx), key)
		return nil
	}

	expNotificationChatIds, err := s.notificationSettings.GetChatsForExpNotification(ctx, req.ExperimentID)
	if err != nil {
		return nil
	}

	text := message.BuildTelegramMessage(&req)

	for _, chatID := range expNotificationChatIds {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		msg.DisableWebPagePreview = true

		_, err = s.bot.Send(msg)
		if err != nil {
			logger.Error.Printf("[%s] TraceId=%s Telegram error: %v", op, observability.GetTraceID(ctx), err)
			return err
		}
	}

	return nil
}
