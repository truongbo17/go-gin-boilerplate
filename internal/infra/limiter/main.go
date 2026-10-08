package limiter

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	redis2 "github.com/truongbo17/go-gin-boilerplate/internal/infra/redis"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

func limitReachedHandler(context *gin.Context) {
	if reset, err := strconv.ParseInt(context.Writer.Header().Get("X-RateLimit-Reset"), 10, 64); err == nil {
		seconds := int64(time.Until(time.Unix(reset, 0)).Seconds())
		if seconds < 1 {
			seconds = 1
		}
		context.Header("Retry-After", strconv.FormatInt(seconds, 10))
	}
	context.AbortWithStatus(http.StatusTooManyRequests)
}

var store limiter.Store

// Limit isolates each policy's counter while sharing the configured backend.
func Limit(policy string, rate limiter.Rate) gin.HandlerFunc {
	if store == nil {
		panic("limiter store is not initialized")
	}
	options := []mgin.Option{
		mgin.WithKeyGetter(func(c *gin.Context) string { return policy + ":" + c.ClientIP() }),
		mgin.WithLimitReachedHandler(limitReachedHandler),
		mgin.WithErrorHandler(func(c *gin.Context, err error) {
			log.Printf("rate limiter store error: %v", err)
			c.AbortWithStatus(http.StatusServiceUnavailable)
		}),
	}

	return mgin.NewMiddleware(limiter.New(store, rate), options...)
}

func InitLimiterStore(storeCache string) {
	options := limiter.StoreOptions{
		Prefix: config.CacheKeyRateLimit,
	}
	if storeCache == config.CacheStoreRedis {
		client := redis2.ClientRedis
		storeRedis, err := redis.NewStoreWithOptions(client, options)
		if err != nil {
			panic(err)
		}

		store = storeRedis
	} else {
		store = memory.NewStoreWithOptions(options)
	}

	fmt.Println("Success init limiter middleware with store " + storeCache)
}
