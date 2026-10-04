package responses

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"time"
)

type UserResponse struct {
	ID        uint         `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Status    enums.Status `json:"status"`
	CreatedAt *time.Time   `json:"created_at"`
	UpdatedAt *time.Time   `json:"updated_at"`
}
