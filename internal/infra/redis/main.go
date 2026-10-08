package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"time"
)

var ClientRedis *redis.Client

func ConnectRedis() error {
	cacheConfig := config.EnvConfig.Cache
	if cacheConfig.CacheStore == config.CacheStoreRedis {
		redisClient := redis.NewClient(&redis.Options{
			Addr:        fmt.Sprintf("%s:%s", cacheConfig.RedisHost, cacheConfig.RedisPort),
			Username:    cacheConfig.RedisUsername,
			Password:    cacheConfig.RedisPassword,
			DialTimeout: 3 * time.Second,
		})

		if err := redisotel.InstrumentTracing(redisClient); err != nil {
			_ = redisClient.Close()
			return fmt.Errorf("instrument Redis: %w", err)
		}
		if err := checkRedisConnection(redisClient); err != nil {
			_ = redisClient.Close()
			return fmt.Errorf("connect Redis: %w", err)
		}
		ClientRedis = redisClient

		fmt.Println("Success connect to Redis.")
	}
	return nil
}

func checkRedisConnection(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}
