package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/truongbo17/go-gin-boilerplate/config"
)

func Open(cacheConfig config.Cache) (*redis.Client, error) {
	if cacheConfig.CacheStore == config.CacheStoreRedis {
		redisClient := redis.NewClient(&redis.Options{
			Addr:        fmt.Sprintf("%s:%s", cacheConfig.RedisHost, cacheConfig.RedisPort),
			Username:    cacheConfig.RedisUsername,
			Password:    cacheConfig.RedisPassword,
			DialTimeout: 3 * time.Second,
		})

		if err := redisotel.InstrumentTracing(redisClient); err != nil {
			_ = redisClient.Close()
			return nil, fmt.Errorf("instrument Redis: %w", err)
		}
		if err := checkRedisConnection(redisClient); err != nil {
			_ = redisClient.Close()
			return nil, fmt.Errorf("connect Redis: %w", err)
		}
		return redisClient, nil
	}
	return nil, nil
}

func checkRedisConnection(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}
