package middlewares

import (
	"github.com/gin-gonic/gin"
	core "github.com/truongbo17/go-gin-boilerplate/internal/infra/limiter"
	"github.com/ulule/limiter/v3"
	"time"
)

func RateAPILimit() gin.HandlerFunc {
	return core.Limit("api", limiter.Rate{Period: time.Minute, Limit: 300})
}

func RateLimitLogin() gin.HandlerFunc {
	return core.Limit("login", limiter.Rate{Period: time.Minute, Limit: 10})
}
