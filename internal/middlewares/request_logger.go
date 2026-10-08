package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/logger"
	"time"
)

func RequestLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		timeNow := time.Now()
		requestId := context.GetString(config.HeaderRequestID)
		clientIp := context.ClientIP()
		userAgent := context.Request.UserAgent()
		method := context.Request.Method
		path := context.Request.URL.Path

		context.Next()

		if path == "/ping" || path == "/ready" {
			return
		}

		logger.LogrusLogger.WithFields(logrus.Fields{
			"request_id": requestId,
			"client_ip":  clientIp,
			"user_agent": userAgent,
			"method":     method,
			"path":       path,
			"status":     context.Writer.Status(),
			"latency_ms": time.Since(timeNow).Milliseconds(),
		}).Info("request completed")
	}
}
