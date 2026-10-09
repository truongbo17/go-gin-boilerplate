package register

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/schedule"
)

// Registry contains the application's background handlers and schedules.
type Registry struct {
	Handlers  map[string]asynq.HandlerFunc
	Schedules []schedule.Job
}

// Module exposes the task handlers owned by one worker feature.
type Module interface {
	Handlers() map[string]asynq.HandlerFunc
}

// New combines enabled worker modules and rejects duplicate task types.
func New(modules ...Module) (Registry, error) {
	registry := Registry{
		Handlers:  make(map[string]asynq.HandlerFunc),
		Schedules: exampleSchedules(),
	}
	for _, module := range modules {
		for taskType, handler := range module.Handlers() {
			if taskType == "" || handler == nil {
				return Registry{}, fmt.Errorf("invalid worker handler for task type %q", taskType)
			}
			if _, exists := registry.Handlers[taskType]; exists {
				return Registry{}, fmt.Errorf("duplicate worker task type %q", taskType)
			}
			registry.Handlers[taskType] = handler
		}
	}
	return registry, nil
}
