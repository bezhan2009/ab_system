package main

import (
	"ab_system/internal/configs"
	"ab_system/internal/repository/redis"
	"ab_system/internal/telegram_bot/bot"
	postgresNotify "ab_system/internal/telegram_bot/repository/postgres"
	"ab_system/pkg/db/postgres"
	redisConn "ab_system/pkg/db/redis"
	"ab_system/pkg/logger"
	"errors"
	"fmt"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(errors.New(fmt.Sprintf("error loading .env file. Error is %s", err)))
	}

	cfg, err := configs.ReadSettings()
	if err != nil {
		panic(errors.New(fmt.Sprintf("error reading settings file. Error is %s", err)))
	}

	err = logger.InitLogger(cfg.LogParams)
	if err != nil {
		panic(errors.New(fmt.Sprintf("error initializing logger. Error is %s", err)))
	}

	dbConn, err := postgres.Connect()
	if err != nil {
		panic(errors.New(fmt.Sprintf("error connecting to database. Error is %s", err)))
	}

	redisClient, err := redisConn.InitializeRedis()
	if err != nil {
		panic(errors.New(fmt.Sprintf("error initializing redis client. Error is %s", err)))
	}

	notificationSettingsRepository := postgresNotify.NewNotificationSettings(dbConn)

	redisRepository := redis.NewRedisRepository(redisClient)

	bot.StartTelegramBot(
		notificationSettingsRepository,
		redisRepository, time.Duration(cfg.NotificationParams.TtlNotificationsMinutes)*time.Minute,
	)
}
