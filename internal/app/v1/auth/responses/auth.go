package responses

import (
	"time"

	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
)

type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Permissions  []string `json:"permissions,omitempty"`
}

type RegisterResponse struct {
	ID          uint `json:"id"`
	EmailQueued bool `json:"email_queued"`
}

type MeResponse struct {
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Status    enums.Status `json:"status"`
	CreateAt  *time.Time   `json:"create_at"`
	UpdatedAt *time.Time   `json:"updated_at"`
}
