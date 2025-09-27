package infra

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisIdempotencyStoreImpl struct{ Client *redis.Client }

func (r *RedisIdempotencyStoreImpl) SetIfNotExists(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	res, err := r.Client.SetNX(ctx, key, value, ttl).Result()
	return res, err
}

func (r *RedisIdempotencyStoreImpl) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func NewRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}
