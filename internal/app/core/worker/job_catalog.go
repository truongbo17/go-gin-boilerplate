package worker

import "fmt"

// TaskType is the stable name stored with a queued task.
type TaskType string

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"

	TypeWelcomeEmail       TaskType = "auth:welcome_email"
	TypePasswordResetEmail TaskType = "auth:password_reset_email"
)

// QueueName scopes an Asynq queue to an application. Asynq stores Redis keys
// under queue names, so this keeps separate apps from sharing queue state.
func QueueName(appName, queue string) string {
	return fmt.Sprintf("%s_%s", appName, queue)
}

type WelcomeEmailParams struct {
	UserID uint `json:"user_id"`
}

type PasswordResetEmailParams struct {
	Email string `json:"email"`
}
