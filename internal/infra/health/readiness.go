package health

import (
	"context"
	"errors"

	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/redis"
)

func Ready(ctx context.Context) error {
	if database.DB == nil {
		return errors.New("database not initialized")
	}
	pool, err := database.DB.DB()
	if err != nil {
		return err
	}
	if err := pool.PingContext(ctx); err != nil {
		return err
	}
	if config.EnvConfig.Cache.CacheStore == config.CacheStoreRedis {
		if redis.ClientRedis == nil {
			return errors.New("redis not initialized")
		}
		return redis.ClientRedis.Ping(ctx).Err()
	}
	return nil
}
