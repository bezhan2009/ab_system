package server

import (
	httpService "ab_system/internal/app"
	"ab_system/internal/clients/telegram_notifications"
	"ab_system/internal/configs"
	"ab_system/internal/repository/postgres"
	redis2 "ab_system/internal/repository/redis"
	postgresConn "ab_system/pkg/db/postgres"
	redisConn "ab_system/pkg/db/redis"
	"ab_system/pkg/logger"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var mainServer *Server

func ServiceStart(ctx context.Context, cfg configs.Configs, db *gorm.DB, redis *redis.Client) (err error) {
	gin.SetMode(cfg.AppParams.GinMode)
	router := gin.Default()

	usersRepository := postgres.NewUsersRepository(db)
	teamRepository := postgres.NewTeamRepository(db)
	featureFlagRepository := postgres.NewFeatureFlagRepository(db)
	experimentsRepository := postgres.NewExperimentsRepository(db)
	experimentVersions := postgres.NewExperimentVersionsRepository(db)
	variantRepository := postgres.NewVariantRepository(db)
	decisionRepository := postgres.NewDecisionRepository(db)
	approverGroupRepository := postgres.NewApproverGroupRepository(db)
	approvalRepository := postgres.NewApprovalRepository(db)
	eventRepository := postgres.NewEventRepository(db)
	eventTypeRepository := postgres.NewEventTypeRepository(db)
	metricRepository := postgres.NewMetricRepository(db)
	expMetricRepository := postgres.NewExperimentMetricRepository(db)
	guardrailTriggerRepository := postgres.NewGuardrailTriggerRepository(db)
	notificationSettings := postgres.NewNotificationSettingsRepository(db)

	redisRepository := redis2.NewRedisRepository(redis)

	telegramClient, err := telegram_notifications.NewNotifyClient(
		cfg.Clients.Telegram.Address,
	)
	if err != nil {
		log.Fatalf("Failed to connect to telegram bot: %v", err)
	}

	httpServer := httpService.NewServer(
		router,
		ctx,
		cfg,
		usersRepository,
		featureFlagRepository,
		experimentsRepository,
		experimentVersions,
		variantRepository,
		decisionRepository,
		teamRepository,
		approverGroupRepository,
		approvalRepository,
		eventRepository,
		eventTypeRepository,
		metricRepository,
		expMetricRepository,
		guardrailTriggerRepository,
		notificationSettings,

		redisRepository,

		telegramClient,
	)

	mainServer = new(Server)
	go func() {
		if err := mainServer.Run(cfg.AppParams.PortRun, cfg.ServerParams, httpServer); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error while starting HTTP Service: %s", err)
		}
	}()

	return nil
}

func ServiceShutdown(db *gorm.DB) {
	fmt.Printf("\n%s\n", "Start of service termination")

	// Закрытие соединения с БД
	err := postgresConn.CloseDBConn(db)
	if err != nil {
		strErr := fmt.Sprintf("Error closing database connection: %s", err.Error())
		fmt.Println(strErr)
		logger.Error.Println(strErr)
	}

	// Закрытие соединения с Redis
	err = redisConn.CloseRedisConnection()
	if err != nil {
		strErr := fmt.Sprintf("Error closing redis connection: %s", err.Error())
		fmt.Println(strErr)
		logger.Error.Println(strErr)
	}

	// Корректное завершение HTTP-сервера
	if err = mainServer.Shutdown(context.Background()); err != nil {
		strErr := fmt.Sprintf("Error shutting down server: %s", err.Error())
		fmt.Println(strErr)
		logger.Error.Println(strErr)
	} else {
		strSuccess := "HTTP-service termination successfully"
		fmt.Println(strSuccess)
		logger.Info.Println(strSuccess)
	}
}
