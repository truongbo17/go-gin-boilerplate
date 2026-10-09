package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/requests"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares"
	"github.com/ulule/limiter/v3"
)

func (module *Module) RegisterRoutes(r *gin.RouterGroup) {
	module.registerAuthRoutes(r)
	module.registerRBACRoutes(r)
}

func (module *Module) rateLimit(policy string, perMinute int64) gin.HandlerFunc {
	return module.limiter.Limit(policy, limiter.Rate{Period: time.Minute, Limit: perMinute})
}

func (module *Module) registerAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group(config.PathV1).Group("auth")
	jwt := middlewares.JWTMiddleware(module.authService)

	auth.POST("register",
		module.rateLimit("auth:register", 5),
		requests.RegisterValidator(),
		module.authController.Register,
	)
	auth.POST("forgot-password",
		module.rateLimit("auth:forgot-password", 5),
		requests.ForgotPasswordValidator(),
		module.authController.ForgotPassword,
	)
	auth.POST("reset-password",
		module.rateLimit("auth:reset-password", 10),
		requests.ResetPasswordValidator(),
		module.authController.ResetPassword,
	)

	auth.POST("login",
		module.rateLimit("auth:login", 10),
		requests.LoginValidator(),
		module.authController.Login,
	)
	auth.GET("me",
		module.rateLimit("auth:me", 300),
		jwt,
		module.authController.Me,
	)
	auth.POST("logout",
		module.rateLimit("auth:logout", 300),
		jwt,
		module.authController.Logout,
	)
	auth.POST("change-pass",
		module.rateLimit("auth:change-password", 300),
		jwt,
		requests.ChangePasswordValidator(),
		module.authController.ChangePass,
	)
}

func (module *Module) registerRBACRoutes(r *gin.RouterGroup) {
	rbac := r.Group(config.PathV1).Group("rbac")
	jwt := middlewares.JWTMiddleware(module.authService)

	rbac.GET("/roles",
		module.rateLimit("role:list", 300),
		jwt,
		middlewares.CheckPermission("role:index", module.permissionService),
		requests.ListRoleValidator(),
		module.roleController.ListRoles,
	)
	rbac.POST("/roles",
		module.rateLimit("role:create", 300),
		jwt,
		middlewares.CheckPermission("role:create", module.permissionService),
		requests.CreateRoleValidator(),
		module.roleController.CreateRole,
	)
	rbac.PUT("/roles/:id",
		module.rateLimit("role:update", 300),
		jwt,
		middlewares.CheckPermission("role:update", module.permissionService),
		requests.UpdateRoleValidator(),
		module.roleController.UpdateRole,
	)
	rbac.DELETE("/roles/:id",
		module.rateLimit("role:delete", 300),
		jwt,
		middlewares.CheckPermission("role:delete", module.permissionService),
		module.roleController.DeleteRole,
	)
	rbac.POST("/roles/:id/permissions",
		module.rateLimit("role:assign-permissions", 300),
		jwt,
		middlewares.CheckPermission("role:assign_permission", module.permissionService),
		requests.AssignPermissionToRoleValidator(),
		module.roleController.AssignPermissionToRole,
	)
	rbac.GET("/roles/:id/permissions",
		module.rateLimit("role:get-permissions", 300),
		jwt,
		middlewares.CheckPermission("role:assign_permission", module.permissionService),
		module.roleController.GetRolePermissions,
	)

	rbac.GET("/permissions",
		module.rateLimit("permission:list", 300),
		jwt,
		middlewares.CheckPermission("permission:index", module.permissionService),
		requests.ListPermissionValidator(),
		module.permissionController.ListPermissions,
	)
	rbac.POST("/permissions",
		module.rateLimit("permission:create", 300),
		jwt,
		middlewares.CheckPermission("permission:create", module.permissionService),
		requests.CreatePermissionValidator(),
		module.permissionController.CreatePermission,
	)
	rbac.PUT("/permissions/:id",
		module.rateLimit("permission:update", 300),
		jwt,
		middlewares.CheckPermission("permission:update", module.permissionService),
		requests.UpdatePermissionValidator(),
		module.permissionController.UpdatePermission,
	)
	rbac.DELETE("/permissions/:id",
		module.rateLimit("permission:delete", 300),
		jwt,
		middlewares.CheckPermission("permission:delete", module.permissionService),
		module.permissionController.DeletePermission,
	)

	rbac.GET("/users",
		module.rateLimit("user:list", 300),
		jwt,
		middlewares.CheckPermission("user:index", module.permissionService),
		requests.ListUserValidator(),
		module.userController.ListUsers,
	)
	rbac.POST("/users/:id/roles",
		module.rateLimit("user:assign-role", 300),
		jwt,
		middlewares.CheckPermission("user:assign_role", module.permissionService),
		requests.AssignRoleToUserValidator(),
		module.userController.AssignRoleToUser,
	)
	rbac.GET("/users/:id/roles",
		module.rateLimit("user:get-roles", 300),
		jwt,
		middlewares.CheckPermission("user:assign_role", module.permissionService),
		module.userController.GetUserRoles,
	)
}
