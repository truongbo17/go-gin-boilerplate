package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return val, err
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) Increment(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	timeValue, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	err = r.client.Expire(ctx, key, ttl).Err()
	if err != nil {
		return 0, err
	}

	return timeValue, nil
}
