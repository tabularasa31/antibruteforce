package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/tabularasa31/antibruteforce/config"

	"time"
)

type Redis struct {
	client *redis.Client
}

func (r Redis) Set(ctx context.Context, key, value string, expire time.Duration) error {
	return r.client.Set(ctx, key, value, expire).Err()
}

func (r Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r Redis) Inc(ctx context.Context, key string) error {
	return r.client.Incr(ctx, key).Err()
}

func (r Redis) Close() error {
	return r.client.Close()
}

func NewRedis(cfg *config.Config) *Redis {

	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	}

	client := redis.NewClient(opt)

	r := &Redis{
		client: client,
	}

	return r
}
