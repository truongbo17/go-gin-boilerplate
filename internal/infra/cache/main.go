package cache

import (
	"context"
	"time"

	redisclient "github.com/redis/go-redis/v9"
	"github.com/truongbo17/go-gin-boilerplate/config"
)

type ICache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	Delete(ctx context.Context, key string) error
}

var _ ICache = (*Redis)(nil)
var _ ICache = (*Local)(nil)

func NewStore(cacheConfig config.Cache, client *redisclient.Client) ICache {
	if cacheConfig.CacheStore == config.CacheStoreRedis {
		return NewRedis(client)
	}
	return NewLocal(5*time.Minute, 10*time.Minute)
}
