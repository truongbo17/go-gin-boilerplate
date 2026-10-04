package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	"net/http"
	"strings"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(config.HeaderAuth)
		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], config.TokenType) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		auth := services.NewAuthService()
		userID, err := auth.VerifyToken(c.Request.Context(), parts[1], enums.TokenTypeAccess)
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		user, err := auth.GetUserById(c.Request.Context(), userID)
		if err != nil || user == nil || user.Status != enums.StatusActive {
			_ = c.Error(errors.New("user unavailable"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
