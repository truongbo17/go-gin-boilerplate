package worker

// TaskType is the stable name stored with a queued task.
type TaskType string

const (
	TypeWelcomeEmail       TaskType = "auth:welcome_email"
	TypePasswordResetEmail TaskType = "auth:password_reset_email"
)

type WelcomeEmailParams struct {
	UserID uint `json:"user_id"`
}

type PasswordResetEmailParams struct {
	Email string `json:"email"`
}
