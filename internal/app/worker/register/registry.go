package register

import (
	"github.com/hibiken/asynq"
	coreworker "github.com/truongbo17/go-gin-boilerplate/internal/app/core/worker"
	authworker "github.com/truongbo17/go-gin-boilerplate/internal/app/worker/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/schedule"
)

// Registry contains the application's background handlers and schedules.
type Registry struct {
	Handlers  map[string]asynq.HandlerFunc
	Schedules []schedule.Job
}

type Handlers struct {
	Auth *authworker.Handler
}

// New registers the handlers and schedules supplied for one worker process.
func New(handlers Handlers) Registry {
	registry := Registry{
		Handlers:  make(map[string]asynq.HandlerFunc),
		Schedules: exampleSchedules(),
	}
	if handlers.Auth != nil {
		registry.Handlers[string(coreworker.TypeWelcomeEmail)] = handlers.Auth.Welcome
		registry.Handlers[string(coreworker.TypePasswordResetEmail)] = handlers.Auth.PasswordReset
	}
	return registry
}
