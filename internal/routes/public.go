package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
)

// LoadPublic registers public routes.
// @Summary Ping endpoint
// @Description Reports that the HTTP process is running.
// @Tags Health
// @Produce json
// @Success 200 {object} response.BaseResponse
// @Router /ping [get]
func LoadPublic(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.BaseResponse{
			Status:     true,
			StatusCode: http.StatusOK,
			RequestId:  c.GetString(config.HeaderRequestID),
		})
	})
}
