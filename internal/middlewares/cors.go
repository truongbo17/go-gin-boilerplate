package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"strings"
)

func Cors() gin.HandlerFunc {
	patternAllowOrigin := config.EnvConfig.Cors.AllowOrigin
	allowOrigin := []string{"http://localhost:3000"}

	if patternAllowOrigin != "" {
		allowOrigin = nil
		for _, origin := range strings.Split(patternAllowOrigin, ",") {
			if origin = strings.TrimSpace(origin); origin != "" {
				allowOrigin = append(allowOrigin, origin)
			}
		}
		if len(allowOrigin) == 0 {
			allowOrigin = []string{"http://localhost:3000"}
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigin,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Language", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset", "Retry-After"},
		AllowCredentials: false,
	})
}
