package health

import (
	"context"
	"errors"

	redisclient "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Checker struct {
	DB         *gorm.DB
	Redis      *redisclient.Client
	CheckRedis bool
}

func (checker Checker) Ready(ctx context.Context) error {
	if checker.DB == nil {
		return errors.New("database not initialized")
	}
	pool, err := checker.DB.DB()
	if err != nil {
		return err
	}
	if err := pool.PingContext(ctx); err != nil {
		return err
	}
	if checker.CheckRedis {
		if checker.Redis == nil {
			return errors.New("redis not initialized")
		}
		return checker.Redis.Ping(ctx).Err()
	}
	return nil
}
