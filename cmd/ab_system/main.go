package main

import (
	"ab_system/internal/configs"
	"ab_system/internal/server"
	"ab_system/pkg/db/postgres"
	redisConn "ab_system/pkg/db/redis"
	"ab_system/pkg/logger"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

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

	err = postgres.Migrate(dbConn)
	if err != nil {
		panic(errors.New(fmt.Sprintf("error migrating database. Error is %s", err)))
	}

	redisClient, err := redisConn.InitializeRedis()
	if err != nil {
		panic(errors.New(fmt.Sprintf("error initializing redis client. Error is %s", err)))
	}

	ctx, cancel := context.WithCancel(context.Background())

	err = server.ServiceStart(ctx, cfg, dbConn, redisClient)
	if err != nil {
		panic(errors.New(fmt.Sprintf("error starting server. Error is %s", err)))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	cancel()

	server.ServiceShutdown(dbConn)
	fmt.Println("Goodbye! :)")
}
