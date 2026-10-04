package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/health"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
	"net/http"
	"time"
)

// LoadPublic sets up public routes.
// @Summary Ping endpoint
// @Description Responds with "pong" and the request ID.
// @Tags Public APi
// @Accept json
// @Produce plain
// @Success 200 {string} string "pong: <x-request-id>"
// @Router /ping [get]
func LoadPublic(r *gin.Engine) *gin.RouterGroup {
	r.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := health.Ready(ctx); err != nil {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}
		c.Status(http.StatusOK)
	})
	public := r.Group("/")
	public.Use(middlewares.RateLimitPublic())
	{
		public.GET("ping", func(context *gin.Context) {
			context.JSON(http.StatusOK, response.BaseResponse{
				Status:     true,
				StatusCode: http.StatusOK,
				RequestId:  context.GetString(config.HeaderRequestID),
			})
		})
	}
	return public
}
