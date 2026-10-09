package middlewares

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckPermission(permission string, permissionService services.PermissionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		ctx := c.Request.Context()
		hasPermission, err := permissionService.CheckPermission(ctx, user.ID, permission)
		if err != nil || !hasPermission {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
