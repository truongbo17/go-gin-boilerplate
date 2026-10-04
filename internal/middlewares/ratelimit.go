package middlewares

import (
	"github.com/gin-gonic/gin"
	core "github.com/truongbo17/go-gin-boilerplate/internal/infra/limiter"
	"github.com/ulule/limiter/v3"
	"time"
)

func RateInternalLimit() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  1000,
	}

	return core.Limit(rate)
}

func RateLimitPublic() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}

	return core.Limit(rate)
}

func RateLimitLogin() gin.HandlerFunc {
	return core.Limit(limiter.Rate{Period: 1 * time.Minute, Limit: 10})
}
