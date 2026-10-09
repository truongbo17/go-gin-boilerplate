package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"time"
)

func RequestLoggerWith(log *logrus.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		timeNow := time.Now()
		requestId := context.GetString(config.HeaderRequestID)
		clientIp := context.ClientIP()
		userAgent := context.Request.UserAgent()
		if len(userAgent) > 256 {
			userAgent = userAgent[:256]
		}
		method := context.Request.Method
		path := context.Request.URL.Path
		if len(path) > 2048 {
			path = path[:2048]
		}

		context.Next()

		if path == "/ping" {
			return
		}

		fields := logrus.Fields{
			"request_id": requestId,
			"client_ip":  clientIp,
			"user_agent": userAgent,
			"method":     method,
			"path":       path,
			"status":     context.Writer.Status(),
			"latency_ms": time.Since(timeNow).Milliseconds(),
		}
		if len(context.Errors) > 0 {
			fields["error"] = context.Errors.Last().Error()
		}
		log.WithFields(fields).Info("request completed")
	}
}
