package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/cache"
)

func TestLogoutRevokesAccessToken(t *testing.T) {
	config.EnvConfig = &config.Config{Auth: config.Auth{JWTSecretKey: strings.Repeat("x", 32), JWTAccessExpirationMinutes: 10}}
	cache.Cache = cache.NewLocal(time.Minute, time.Minute)
	service := AuthService{}
	user := &models.User{}
	user.ID = 42
	token, err := service.GenerateToken(context.Background(), enums.TokenTypeAccess, user)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.VerifyToken(context.Background(), token, enums.TokenTypeAccess); err != nil {
		t.Fatal(err)
	}
	if failure := service.Logout(context.Background(), types.LogoutInput{User: user, Token: "Bearer " + token}); failure != nil {
		t.Fatal(failure)
	}
	if _, err := service.VerifyToken(context.Background(), token, enums.TokenTypeAccess); err == nil {
		t.Fatal("revoked token was accepted")
	}
}

func TestPasswordChangeInvalidatesPreviousTokenVersion(t *testing.T) {
	config.EnvConfig = &config.Config{Auth: config.Auth{JWTSecretKey: strings.Repeat("x", 32), JWTAccessExpirationMinutes: 10}}
	cache.Cache = cache.NewLocal(time.Minute, time.Minute)
	service := AuthService{}
	user := &models.User{Password: "previous-password-hash"}
	user.ID = 42
	token, err := service.GenerateToken(context.Background(), enums.TokenTypeAccess, user)
	if err != nil {
		t.Fatal(err)
	}
	_, claims, err := service.VerifyTokenClaims(context.Background(), token, enums.TokenTypeAccess)
	if err != nil || claims.PasswordVersion != PasswordVersionForUser(user) {
		t.Fatalf("new token has invalid password version: %v", err)
	}
	user.Password = "new-password-hash"
	if claims.PasswordVersion == PasswordVersionForUser(user) {
		t.Fatal("old token remained valid after password change")
	}
}
