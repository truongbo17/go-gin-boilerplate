package auth

import (
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/controllers"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/requests"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadAuthV1(r *gin.RouterGroup) {
	handler := controllers.NewAuthController()

	auth := r.Group(config.PathV1).Group("auth")
	{
		auth.GET("me",
			middlewares.JWTMiddleware(),
			handler.Me,
		)
		auth.POST("login",
			middlewares.RateLimitLogin(),
			requests.LoginValidator(),
			handler.Login,
		)
		auth.POST(
			"logout",
			middlewares.JWTMiddleware(),
			handler.Logout,
		)
		auth.POST(
			"/change-pass",
			requests.ChangePasswordValidator(),
			middlewares.JWTMiddleware(),
			handler.ChangePass,
		)
	}
}

func LoadRBACV1(r *gin.RouterGroup) {
	roleController := controllers.NewRoleController()
	permissionController := controllers.NewPermissionController()
	userController := controllers.NewUserController()

	rbac := r.Group(config.PathV1).Group("rbac", middlewares.JWTMiddleware())
	{
		rbac.GET("/roles", middlewares.CheckPermission("role:index"), requests.ListRoleValidator(), roleController.ListRoles)
		rbac.POST("/roles", middlewares.CheckPermission("role:create"), requests.CreateRoleValidator(), roleController.CreateRole)
		rbac.PUT("/roles/:id", middlewares.CheckPermission("role:update"), requests.UpdateRoleValidator(), roleController.UpdateRole)
		rbac.DELETE("/roles/:id", middlewares.CheckPermission("role:delete"), roleController.DeleteRole)
		rbac.POST("/roles/:id/permissions", middlewares.CheckPermission("role:assign_permission"), requests.AssignPermissionToRoleValidator(), roleController.AssignPermissionToRole)

		rbac.GET("/permissions", middlewares.CheckPermission("permission:index"), requests.ListPermissionValidator(), permissionController.ListPermissions)
		rbac.POST("/permissions", middlewares.CheckPermission("permission:create"), requests.CreatePermissionValidator(), permissionController.CreatePermission)
		rbac.PUT("/permissions/:id", middlewares.CheckPermission("permission:update"), requests.UpdatePermissionValidator(), permissionController.UpdatePermission)
		rbac.DELETE("/permissions/:id", middlewares.CheckPermission("permission:delete"), permissionController.DeletePermission)

		rbac.GET("/users", middlewares.CheckPermission("user:index"), requests.ListUserValidator(), userController.ListUsers)
		rbac.POST("/users/:id/roles", middlewares.CheckPermission("user:assign_role"), requests.AssignRoleToUserValidator(), roleController.AssignRoleToUser)
		rbac.GET("/users/:id/roles", middlewares.CheckPermission("user:assign_role"), roleController.GetUserRoles)
		rbac.GET("/roles/:id/permissions", middlewares.CheckPermission("role:assign_permission"), roleController.GetRolePermissions)
	}
}
