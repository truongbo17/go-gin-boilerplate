package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000")
		}
		c.Next()
	}
}

func LimitRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > 1<<20 {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	}
}
