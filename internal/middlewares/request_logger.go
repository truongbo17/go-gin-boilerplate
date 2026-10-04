package middlewares

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/logger"
	"time"
)

type RequestLog struct {
	RequestID string  `json:"request_id"`
	ClientIP  string  `json:"client_ip"`
	UserAgent string  `json:"user_agent"`
	Method    string  `json:"method"`
	Path      string  `json:"path"`
	Latency   float64 `json:"latency"`
	Timestamp string  `json:"timestamp"`
}

func RequestLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		timeNow := time.Now()
		requestId := context.GetString(config.HeaderRequestID)
		clientIp := context.ClientIP()
		userAgent := context.Request.UserAgent()
		method := context.Request.Method
		path := context.Request.URL.Path

		context.Next()

		if path == "/ping" {
			return
		}

		latency := time.Since(timeNow).Seconds()

		logData := RequestLog{
			RequestID: requestId,
			ClientIP:  clientIp,
			UserAgent: userAgent,
			Method:    method,
			Path:      path,
			Latency:   latency,
			Timestamp: timeNow.Format(time.DateTime),
		}

		log := logger.LogrusLogger

		logDataJSON, err := json.Marshal(logData)
		if err != nil {
			logrus.Error("Failed to marshal log data to JSON: ", err)
		}

		log.WithFields(logrus.Fields{
			"log": string(logDataJSON),
		}).Infoln("request logger")
	}
}
