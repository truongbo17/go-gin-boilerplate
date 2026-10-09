package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

func CorsWith(patternAllowOrigin string) gin.HandlerFunc {
	allowOrigin := []string{"http://localhost:3000"}

	if patternAllowOrigin != "" {
		allowOrigin = nil
		for origin := range strings.SplitSeq(patternAllowOrigin, ",") {
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
