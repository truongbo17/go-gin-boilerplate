package types

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Username        string          `json:"username"`
	Email           string          `json:"email"`
	Type            enums.TokenType `json:"type"`
	PasswordVersion string          `json:"password_version"`
	Permissions     []string        `json:"permissions,omitempty"`
}
