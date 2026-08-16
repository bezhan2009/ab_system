package handlers

import (
	telegramv1 "ab_system/internal/telegram_bot/gen/telegram_notifications/v1"
	"ab_system/internal/telegram_bot/service"
	"ab_system/pkg/logger"
	notifydto "ab_system/pkg/notifications/dto"
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type NotifyGRPCServer struct {
	telegramv1.UnimplementedTelegramBotServiceServer
	apiSecret string
	notifier  *service.TelegramNotifier
	bot       *tgbotapi.BotAPI
}

func NewNotifyGRPCServer(apiSecret string, notifier *service.TelegramNotifier, bot *tgbotapi.BotAPI) *NotifyGRPCServer {
	return &NotifyGRPCServer{
		apiSecret: apiSecret,
		notifier:  notifier,
		bot:       bot,
	}
}

func (s *NotifyGRPCServer) Notify(ctx context.Context, r *telegramv1.NotifyRequest) (*telegramv1.NotifyResponse, error) {
	if s.apiSecret != "" {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authVals := md.Get("authorization")

		if len(authVals) == 0 || authVals[0] != "Bearer "+s.apiSecret {
			logger.Warn.Println("Unauthorized /notify request")
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}
	}

	if r.EventType == "" {
		return nil, status.Error(codes.InvalidArgument, "event_type is required")
	}

	req := notifydto.NotifyRequest{
		EventType:    r.EventType,
		ExperimentID: r.ExperimentId,
		Experiment:   r.Experiment,
		FlagKey:      r.FlagKey,
		UserID:       r.UserId,
		Status:       r.Status,
		Metric:       r.Metric,
		Threshold:    r.Threshold,
		Value:        r.Value,
		Message:      r.Message,
		ReportURL:    r.ReportUrl,
		Timestamp:    r.Timestamp,
	}

	if err := s.notifier.Send(ctx, req); err != nil {
		return nil, status.Error(codes.Internal, "something went wrong")
	}

	logger.Info.Printf("Sent notification: %s / %s", req.EventType, req.ExperimentID)
	return &telegramv1.NotifyResponse{Message: "successfully sent notification"}, nil
}

func (s *NotifyGRPCServer) Health(ctx context.Context, r *telegramv1.HealthRequest) (*telegramv1.HealthResponse, error) {
	return &telegramv1.HealthResponse{Status: "ok"}, nil
}

func (s *NotifyGRPCServer) Ready(ctx context.Context, r *telegramv1.ReadyRequest) (*telegramv1.ReadyResponse, error) {
	_, err := s.bot.GetMe()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "not ready: %v", err)
	}

	return &telegramv1.ReadyResponse{Status: "ready"}, nil
}
