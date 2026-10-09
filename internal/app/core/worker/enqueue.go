package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// Dispatcher is the port used by services that submit background work.
type Dispatcher interface {
	Dispatch(context.Context, TaskType, any) error
}

type JobClient interface {
	EnqueueContext(context.Context, *asynq.Task, ...asynq.Option) (*asynq.TaskInfo, error)
}

type JobDispatcher struct {
	client JobClient
}

func New(client JobClient) *JobDispatcher {
	return &JobDispatcher{client: client}
}

func (dispatcher *JobDispatcher) Dispatch(ctx context.Context, taskType TaskType, params any) error {
	var options []asynq.Option
	switch taskType {
	case TypeWelcomeEmail:
		value, ok := params.(WelcomeEmailParams)
		if !ok || value.UserID == 0 {
			return fmt.Errorf("invalid parameters for task %q", taskType)
		}
		options = []asynq.Option{asynq.Queue("default"), asynq.MaxRetry(5)}
	case TypePasswordResetEmail:
		value, ok := params.(PasswordResetEmailParams)
		if !ok || value.Email == "" {
			return fmt.Errorf("invalid parameters for task %q", taskType)
		}
		options = []asynq.Option{asynq.Queue("critical"), asynq.MaxRetry(5), asynq.Unique(time.Minute)}
	default:
		return fmt.Errorf("unknown task type %q", taskType)
	}
	payload, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal task %q: %w", taskType, err)
	}
	_, err = dispatcher.client.EnqueueContext(ctx, asynq.NewTask(string(taskType), payload), options...)
	if taskType == TypePasswordResetEmail && errors.Is(err, asynq.ErrDuplicateTask) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("enqueue task %q: %w", taskType, err)
	}
	return nil
}
