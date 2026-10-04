package limiter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	redis2 "github.com/truongbo17/go-gin-boilerplate/internal/infra/redis"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"net/http"
	"strconv"
)

func limitReachedHandler(context *gin.Context) {
	context.AbortWithStatus(http.StatusTooManyRequests)
}

var store limiter.Store

func keyGetterUserID(c *gin.Context) string {
	return strconv.FormatInt(c.MustGet("userId").(int64), 16)
}

func keyGetterIP(context *gin.Context) string {
	return context.ClientIP()
}

func Limit(rate limiter.Rate) gin.HandlerFunc {
	options := []mgin.Option{
		mgin.WithKeyGetter(keyGetterIP),
		mgin.WithLimitReachedHandler(limitReachedHandler),
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
