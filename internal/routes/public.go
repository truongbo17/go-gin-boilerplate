package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
	"net/http"
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
