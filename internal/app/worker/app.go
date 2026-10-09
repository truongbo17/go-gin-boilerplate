package worker

import (
	redisclient "github.com/redis/go-redis/v9"
	"github.com/truongbo17/go-gin-boilerplate/config"
	authworker "github.com/truongbo17/go-gin-boilerplate/internal/app/worker/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/worker/register"
	"gorm.io/gorm"
)

// Dependencies are process resources shared by worker modules.
type Dependencies struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redisclient.Client
}

// New builds enabled worker modules and registers their handlers and schedules.
func New(deps Dependencies) (register.Registry, error) {
	var handlers register.Handlers
	if deps.Config.Mail.Enabled {
		authHandler, err := authworker.NewHandler(authworker.Dependencies{
			Mail:         deps.Config.Mail,
			JWTSecretKey: deps.Config.Auth.JWTSecretKey,
			DB:           deps.DB,
			Redis:        deps.Redis,
		})
		if err != nil {
			return register.Registry{}, err
		}
		handlers.Auth = authHandler
	}
	return register.New(handlers), nil
}
