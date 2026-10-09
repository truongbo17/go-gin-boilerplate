package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	coreworker "github.com/truongbo17/go-gin-boilerplate/internal/app/core/worker"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/controllers"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares/limiter"
	authrepository "github.com/truongbo17/go-gin-boilerplate/internal/repository/auth"
	"gorm.io/gorm"
)

// Module owns the auth feature's repository, service, controller, and route wiring.
type Module struct {
	authService       services.AuthService
	permissionService services.PermissionService
	limiter           *limiter.Limiter

	authController       *controllers.AuthController
	userController       *controllers.UserController
	roleController       *controllers.RoleController
	permissionController *controllers.PermissionController
}

func (module *Module) AuthMiddleware() gin.HandlerFunc {
	return middlewares.JWTMiddleware(module.authService)
}

type Dependencies struct {
	DB             *gorm.DB
	TokenBlacklist services.TokenBlacklist
	Auth           config.Auth
	Limiter        *limiter.Limiter
	Jobs           coreworker.Dispatcher
}

func NewModule(params Dependencies) (*Module, error) {
	if params.DB == nil || params.TokenBlacklist == nil || params.Limiter == nil {
		return nil, errors.New("auth module requires DB, token blacklist, and limiter")
	}
	userRepository := authrepository.NewUserRepository(params.DB)
	roleRepository := authrepository.NewRoleRepository(params.DB)
	permissionRepository := authrepository.NewPermissionRepository(params.DB)

	authService := services.NewAuthService(userRepository, params.TokenBlacklist, services.AuthSettings{
		JWTSecretKey:               params.Auth.JWTSecretKey,
		JWTAccessExpirationMinutes: params.Auth.JWTAccessExpirationMinutes,
		JWTRefreshExpirationDays:   params.Auth.JWTRefreshExpirationDays,
	})
	authService.Jobs = params.Jobs
	userService := services.NewUserService(userRepository)
	roleService := services.NewRoleService(roleRepository)
	permissionService := services.NewPermissionService(permissionRepository)

	return &Module{
		authService:          authService,
		permissionService:    permissionService,
		limiter:              params.Limiter,
		authController:       controllers.NewAuthController(authService),
		userController:       controllers.NewUserController(userService),
		roleController:       controllers.NewRoleController(roleService, permissionService),
		permissionController: controllers.NewPermissionController(permissionService),
	}, nil
}
