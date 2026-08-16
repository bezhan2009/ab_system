package telegram_notifications

import (
	telegramv1 "ab_system/internal/telegram_bot/gen/telegram_notifications/v1"
	"ab_system/pkg/logger"
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type NotifyClient struct {
	grpc telegramv1.TelegramBotServiceClient
}

func NewNotifyClient(address string) (*NotifyClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("telegram grpc dial: %w", err)
	}

	return &NotifyClient{
		grpc: telegramv1.NewTelegramBotServiceClient(conn),
	}, nil
}

func (c *NotifyClient) Notify(ctx context.Context, req *telegramv1.NotifyRequest) (err error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+os.Getenv("RANDOM_SECRET"))

	_, err = c.grpc.Notify(ctx, req)
	if err != nil {
		logger.Error.Printf("[clients.telegram_notifications.Notify] Error: %v", err)

		return err
	}

	return nil
}

func (c *NotifyClient) Ready(ctx context.Context) (err error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer ")

	_, err = c.grpc.Ready(ctx, &telegramv1.ReadyRequest{})
	if err != nil {
		logger.Error.Printf("[clients.telegram_notifications.Notify] Error: %v", err)

		return err
	}

	return nil
}
