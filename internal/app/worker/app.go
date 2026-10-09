package worker

import (
	"fmt"

	redisclient "github.com/redis/go-redis/v9"
	"github.com/truongbo17/go-gin-boilerplate/config"
	authworker "github.com/truongbo17/go-gin-boilerplate/internal/app/worker/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/worker/register"
	"gorm.io/gorm"
)

// New lists enabled worker features and builds their handler registry.
func New(cfg config.Config, db *gorm.DB, redis *redisclient.Client) (register.Registry, error) {
	var modules []register.Module
	if cfg.Mail.Enabled {
		handler, err := authworker.NewHandler(authworker.Dependencies{
			DB:           db,
			Redis:        redis,
			Mail:         cfg.Mail,
			JWTSecretKey: cfg.Auth.JWTSecretKey,
		})
		if err != nil {
			return register.Registry{}, fmt.Errorf("build auth mail worker: %w", err)
		}
		modules = append(modules, handler)
	}
	return register.New(modules...)
}
