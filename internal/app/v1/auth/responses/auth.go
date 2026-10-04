package responses

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"time"
)

type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Permissions  []string `json:"permissions,omitempty"`
}

type MeResponse struct {
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Status    enums.Status `json:"status"`
	CreateAt  *time.Time   `json:"create_at"`
	UpdatedAt *time.Time   `json:"updated_at"`
}
