package cache

import (
	"context"
	"fmt"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/redis"
	"time"
)

type ICache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
}

var _ ICache = (*Redis)(nil)
var _ ICache = (*Local)(nil)

var Cache ICache

func InitCache() {
	cacheConfig := config.EnvConfig.Cache
	if cacheConfig.CacheStore == config.CacheStoreRedis {
		Cache = NewRedis(redis.ClientRedis)
	} else {
		Cache = NewLocal(5*time.Minute, 10*time.Minute)
	}

	fmt.Println("Success init cache with store " + cacheConfig.CacheStore)
}
