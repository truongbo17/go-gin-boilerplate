package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"strings"
)

func Cors() gin.HandlerFunc {
	EnvConfig := config.EnvConfig
	patternAllowOrigin := EnvConfig.Cors.AllowOrigin
	allowOrigin := []string{"http://localhost:3000"}

	if patternAllowOrigin != "" {
		allowOrigin = strings.Split(EnvConfig.Cors.AllowOrigin, ",")
		if allowOrigin == nil || len(allowOrigin) == 0 {
			allowOrigin = []string{"http://localhost:3000"}
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigin,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Language", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	})
}
