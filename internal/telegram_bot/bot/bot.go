package bot

import (
	"ab_system/internal/domain/redis"
	telegramv1 "ab_system/internal/telegram_bot/gen/telegram_notifications/v1"
	"ab_system/internal/telegram_bot/handlers"
	interceptors "ab_system/internal/telegram_bot/interceptrors"
	"ab_system/internal/telegram_bot/repository/postgres"
	"ab_system/internal/telegram_bot/service"
	"ab_system/pkg/logger"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartTelegramBot(
	notificationSettings *postgres.NotificationSettingsRepository,
	redisCacheRepository redis.RedisCacheRepository,
	ttlNotify time.Duration) {

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	//chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	apiSecret := os.Getenv("RANDOM_SECRET")
	host := os.Getenv("TELEGRAM_BOT_HOST")
	port := os.Getenv("TELEGRAM_BOT_PORT")

	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "50051"
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Failed to initialize bot: %v", err)))
	}

	me, err := bot.GetMe()
	if err != nil {
		panic(errors.New(fmt.Sprintf("Failed to get bot info: %v", err)))
	}

	logger.Info.Printf("Bot started: @%s", me.UserName)

	notifier := service.NewTelegramNotifier(
		bot,
		notificationSettings,
		redisCacheRepository,
		ttlNotify,
	)
	grpcServer := handlers.NewNotifyGRPCServer(apiSecret, notifier, bot)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.TraceUnaryInterceptor,
		),
	)

	telegramv1.RegisterTelegramBotServiceServer(srv, grpcServer)
	reflection.Register(srv)

	addr := host + ":" + port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Failed to listen: %v", err)))
	}

	logger.Info.Printf("gRPC server listening on %s", addr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info.Println("Received shutdown signal, gracefully stopping server...")
		srv.GracefulStop()
		logger.Info.Println("Server stopped")
	}()

	if err = srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		panic(errors.New(fmt.Sprintf("Failed to serve: %v", err)))
	}
}
