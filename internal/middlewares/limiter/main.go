package limiter

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	redisclient "github.com/redis/go-redis/v9"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

func limitReachedHandler(context *gin.Context) {
	if reset, err := strconv.ParseInt(context.Writer.Header().Get("X-RateLimit-Reset"), 10, 64); err == nil {
		remaining := time.Until(time.Unix(reset, 0))
		seconds := max(int64((remaining+time.Second-1)/time.Second), 1)
		context.Header("Retry-After", strconv.FormatInt(seconds, 10))
	}
	context.AbortWithStatus(http.StatusTooManyRequests)
}

type Limiter struct {
	store limiter.Store
}

func New(storeCache string, client *redisclient.Client) (*Limiter, error) {
	options := limiter.StoreOptions{Prefix: config.CacheKeyRateLimit}
	if storeCache == config.CacheStoreRedis {
		if client == nil {
			return nil, fmt.Errorf("Redis limiter requires a Redis client")
		}
		redisStore, err := redis.NewStoreWithOptions(client, options)
		if err != nil {
			return nil, fmt.Errorf("create Redis limiter store: %w", err)
		}
		return &Limiter{store: redisStore}, nil
	}
	return &Limiter{store: memory.NewStoreWithOptions(options)}, nil
}

// Limit isolates each policy's counter while sharing the configured backend.
func (l *Limiter) Limit(policy string, rate limiter.Rate) gin.HandlerFunc {
	options := []mgin.Option{
		mgin.WithKeyGetter(func(c *gin.Context) string { return policy + ":" + c.ClientIP() }),
		mgin.WithLimitReachedHandler(limitReachedHandler),
		mgin.WithErrorHandler(func(c *gin.Context, err error) {
			_ = c.Error(err)
			c.AbortWithStatus(http.StatusServiceUnavailable)
		}),
	}
	return mgin.NewMiddleware(limiter.New(l.store, rate), options...)
}
