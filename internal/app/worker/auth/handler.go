package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	coreworker "github.com/truongbo17/go-gin-boilerplate/internal/app/core/worker"
)

type Mailer interface {
	SendWelcome(context.Context, uint) error
	SendPasswordReset(context.Context, string) error
}

type Handler struct {
	Mail Mailer
}

func (handler *Handler) Handlers() map[string]asynq.HandlerFunc {
	return map[string]asynq.HandlerFunc{
		string(coreworker.TypeWelcomeEmail):       handler.Welcome,
		string(coreworker.TypePasswordResetEmail): handler.PasswordReset,
	}
}

func (handler *Handler) Welcome(ctx context.Context, task *asynq.Task) error {
	var payload coreworker.WelcomeEmailParams
	if err := json.Unmarshal(task.Payload(), &payload); err != nil || payload.UserID == 0 {
		return fmt.Errorf("invalid welcome task: %w", asynq.SkipRetry)
	}
	return handleMailError(handler.Mail.SendWelcome(ctx, payload.UserID))
}

func (handler *Handler) PasswordReset(ctx context.Context, task *asynq.Task) error {
	var payload coreworker.PasswordResetEmailParams
	if err := json.Unmarshal(task.Payload(), &payload); err != nil || payload.Email == "" {
		return fmt.Errorf("invalid password reset task: %w", asynq.SkipRetry)
	}
	return handleMailError(handler.Mail.SendPasswordReset(ctx, payload.Email))
}

func handleMailError(err error) error {
	if errors.Is(err, services.ErrInvalidMailAddress) {
		return fmt.Errorf("invalid auth mail task: %w", asynq.SkipRetry)
	}
	return err
}
