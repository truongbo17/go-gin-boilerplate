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

func ConnectRedis() {
	cacheConfig := config.EnvConfig.Cache
	if cacheConfig.CacheStore == config.CacheStoreRedis {
		redisClient := redis.NewClient(&redis.Options{
			Addr:        fmt.Sprintf("%s:%s", cacheConfig.RedisHost, cacheConfig.RedisPort),
			Username:    cacheConfig.RedisUsername,
			Password:    cacheConfig.RedisPassword,
			DialTimeout: 3 * time.Second,
		})

		if err := redisotel.InstrumentTracing(redisClient); err != nil {
			panic(err)
		}

		ClientRedis = redisClient

		checkRedisConnection(redisClient)

		fmt.Println("Success connect to Redis.")
	}
}

func checkRedisConnection(client *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}
}
