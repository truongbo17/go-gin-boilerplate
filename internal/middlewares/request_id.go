package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil || !span.SpanContext().IsValid() {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

func RequestID() gin.HandlerFunc {
	return func(context *gin.Context) {
		id := GetTraceID(context.Request.Context())
		if id == "" {
			id = uuid.New().String()
			span := trace.SpanFromContext(context.Request.Context())

			if span != nil && span.SpanContext().IsValid() {
				span.SetAttributes(attribute.String("request_id", id))
			}
		}

		context.Set(config.HeaderRequestID, id)
		context.Header(config.HeaderRequestID, id)

		context.Next()
	}
}
