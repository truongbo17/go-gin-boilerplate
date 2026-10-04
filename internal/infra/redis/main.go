package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/truongbo17/go-gin-boilerplate/config"
)

var ClientRedis *redis.Client

func ConnectRedis() {
	cacheConfig := config.EnvConfig.Cache
	if cacheConfig.RedisHost != "" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cacheConfig.RedisHost, cacheConfig.RedisPort),
			Username: cacheConfig.RedisUsername,
			Password: cacheConfig.RedisPassword,
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
	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}
