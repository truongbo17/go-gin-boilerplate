package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	json "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func StartWith(configRedis config.Cache, log *logrus.Logger, handlers map[string]asynq.HandlerFunc) (*asynq.Server, error) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%s", configRedis.RedisHost, configRedis.RedisPort),
			Username: configRedis.RedisUsername,
			Password: configRedis.RedisPassword,
		},
		asynq.Config{
			Concurrency:     20,
			ShutdownTimeout: 25 * time.Second,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			RetryDelayFunc: retryDelay,
			Logger:         log,
			ErrorHandler:   asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) { errorHandler(log, ctx, task, err) }),
			IsFailure:      func(err error) bool { return !isRateLimitError(err) },
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(workerTracingMiddleware)

	for name, handler := range handlers {
		mux.HandleFunc(name, handler)
	}

	if err := srv.Start(mux); err != nil {
		return nil, fmt.Errorf("start worker: %w", err)
	}
	fmt.Println("Success init server asynq queue.")
	return srv, nil
}

func retryDelay(n int, err error, _ *asynq.Task) time.Duration {
	var delayable DelayableError
	if errors.As(err, &delayable) {
		return max(0, delayable.RetryIn())
	}
	if n < 0 {
		n = 0
	}
	if n > 6 {
		n = 6
	}
	return min(time.Duration(5<<n)*time.Second, 5*time.Minute)
}

func workerTracingMiddleware(next asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		var payload map[string]any
		if err := json.Unmarshal(task.Payload(), &payload); err == nil {
			if traceMap, ok := payload["tracer"].(map[string]any); ok {
				carrier := propagation.MapCarrier{}
				for k, v := range traceMap {
					if strVal, ok := v.(string); ok {
						carrier[k] = strVal
					}
				}
				ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
			}
		}

		tracer := otel.Tracer("AsynqWorker")
		ctx, span := tracer.Start(ctx, "Asynq/ProcessTask")
		defer span.End()

		span.SetAttributes(
			attribute.String("task.id", task.ResultWriter().TaskID()),
			attribute.String("task.type", task.Type()),
		)

		err := next.ProcessTask(ctx, task)
		if err != nil {
			span.RecordError(err)
		}
		return err
	})
}

type RateLimitError struct {
	RetryInDuration time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limited (retry in %v)", e.RetryInDuration)
}

func (e *RateLimitError) RetryIn() time.Duration {
	return e.RetryInDuration
}

type DelayableError interface {
	RetryIn() time.Duration
}

func isRateLimitError(err error) bool {
	var de DelayableError
	return errors.As(err, &de)
}

func errorHandler(log *logrus.Logger, ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, hasMaxRetry := asynq.GetMaxRetry(ctx)
	taskID, _ := asynq.GetTaskID(ctx)
	queue, _ := asynq.GetQueueName(ctx)
	fields := logrus.Fields{
		"task_id": taskID, "task_type": task.Type(), "queue": queue,
		"retry_count": retried, "max_retry": maxRetry, "error": err.Error(),
	}
	if errors.Is(err, asynq.RevokeTask) {
		fields["outcome"] = "revoked"
		log.WithFields(fields).Warn("job revoked")
		return
	}
	if errors.Is(err, asynq.SkipRetry) || (hasMaxRetry && retried >= maxRetry) {
		fields["outcome"] = "archived"
		log.WithFields(fields).Error("job failed permanently")
		return
	}
	fields["outcome"] = "retrying"
	if isRateLimitError(err) {
		log.WithFields(fields).Info("job delayed")
		return
	}
	log.WithFields(fields).Warn("job failed; retry scheduled")
}
