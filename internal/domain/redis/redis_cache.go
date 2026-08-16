package redis

import (
	"context"
	"time"
)

type RedisCacheRepository interface {
	SetCache(ctx context.Context, key string, value interface{}) error
	GetCache(ctx context.Context, key string, dest interface{}) (bool, error)
	DeleteCache(ctx context.Context, key string) error
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)
	Exists(ctx context.Context, key string) (bool, error)
}
