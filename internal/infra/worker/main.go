package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	json "github.com/json-iterator/go"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/register"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"time"
)

func InitServer() {
	EnvConfig := config.EnvConfig
	configRedis := EnvConfig.Cache

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%s", configRedis.RedisHost, configRedis.RedisPort),
			Username: configRedis.RedisUsername,
			Password: configRedis.RedisPassword,
		},
		asynq.Config{
			Concurrency: 20,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			RetryDelayFunc: func(n int, err error, task *asynq.Task) time.Duration {
				var de DelayableError
				if errors.As(err, &de) {
					return de.RetryIn()
				}
				return time.Duration(5<<n) * time.Second // Exponential: 5, 10, 15, 20s...
			},
			Logger:       logger.LogrusLogger,
			ErrorHandler: asynq.ErrorHandlerFunc(errorHandler),
			IsFailure:    func(err error) bool { return !isRateLimitError(err) },
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(asynqTracingMiddleware)

	for name, handler := range register.Handlers {
		mux.HandleFunc(name, handler)
	}

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
	fmt.Println("Success init server asynq queue.")
}

func asynqTracingMiddleware(next asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		var payload map[string]interface{}
		if err := json.Unmarshal(task.Payload(), &payload); err == nil {
			if traceMap, ok := payload["tracer"].(map[string]interface{}); ok {
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

func errorHandler(ctx context.Context, task *asynq.Task, err error) {
	if isRateLimitError(err) {
		return
	}

	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry {
		err = fmt.Errorf("retry exhausted for task %q: %w", task.Type(), err)
	}

	logApp := logger.LogrusLogger
	logApp.Errorf("Job error: task=%q, payload=%q error=%v", task.Type(), task.Payload(), err)

	logger.LogrusLogger.Errorf("Job failed: type=%q error=%v", task.Type(), err)
}
