package redis

import (
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/errs"
	"ab_system/pkg/logger"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) SetCache(ctx context.Context, key string, value interface{}) error {
	const op = "repository.redis.SetCache"

	data, err := json.Marshal(value)
	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return err
	}

	err = r.client.Set(
		ctx,
		key,
		data,
		0,
	).Err()

	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return err
	}

	return nil
}

func (r *RedisRepository) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	const op = "repository.redis.SetNX"

	data, err := json.Marshal(value)
	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v", op, observability.GetTraceID(ctx), err)
		return false, err
	}

	success, err := r.client.SetNX(ctx, key, data, ttl).Result()
	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v", op, observability.GetTraceID(ctx), err)
		return false, err
	}

	return success, nil
}

func (r *RedisRepository) GetCache(ctx context.Context, key string, dest interface{}) (bool, error) {
	const op = "repository.redis.GetCache"

	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		logger.Error.Printf("[db.GetCache] Key %s does not exist", key)

		return false, errs.ErrMissingCache
	} else if err != nil {
		logger.Error.Printf("[db.GetCache] Error getting key %s from redis: %v", key, err)

		return false, err
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Key=%s Error: %v",
			op, observability.GetTraceID(ctx), key, err)

		return false, err
	}

	return true, nil
}

func (r *RedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	const op = "repository.redis.Exists"

	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Key=%s Error: %v",
			op, observability.GetTraceID(ctx), key, err)
		return false, err
	}

	return result > 0, nil
}

func (r *RedisRepository) DeleteCache(ctx context.Context, key string) error {
	const op = "repository.redis.DeleteCache"

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Key=%s Error: %v",
			op, observability.GetTraceID(ctx), key, err)

		return err
	}

	return nil
}
