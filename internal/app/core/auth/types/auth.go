package types

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginOutput struct {
	AccessToken string `json:"access_token"`
}

type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterOutput struct {
	User *models.User `json:"user"`
}

type LogoutInput struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}
