package db

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func InitializeRedis() (*redis.Client, error) {
	var addr string

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	if os.Getenv("REDIS_HOST") != "" {
		addr = fmt.Sprintf("%s:%s", redisHost, redisPort)
	} else {
		addr = ":6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword,
		DB:       redisDB,
	})

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return RedisClient, nil
}

func CloseRedisConnection() error {
	err := RedisClient.Close()
	if err != nil {
		return err
	}

	return nil
}
