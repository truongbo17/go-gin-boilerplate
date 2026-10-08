package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/repositories"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/cache"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type AuthService struct {
	UserRepository repositories.UserRepository
}

func NewAuthService() AuthService {
	return AuthService{UserRepository: repositories.NewUserRepository()}
}

func (authService *AuthService) Login(ctx context.Context, input types.LoginInput) (*types.LoginOutput, *core.ErrorReturn) {
	user, err := authService.UserRepository.FindOneByCondition(ctx, models.User{Username: input.Username})
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrAuthLoginFailed,
			Err:       err,
		}
	}
	if user == nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrAuthLoginFailed,
		}
	}
	if user.Status != enums.StatusActive {
		return nil, &core.ErrorReturn{ErrorCode: response.ErrAuthLoginFailed}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrAuthLoginFailed,
		}
	}

	token, err := authService.generateToken(ctx, enums.TokenTypeAccess, user)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrAuthGenerateToken,
			Err:       err,
		}
	}

	return &types.LoginOutput{
		AccessToken: token,
	}, nil
}

func (authService *AuthService) Register(ctx context.Context, input types.RegisterInput) (*types.RegisterOutput, *core.ErrorReturn) {
	exist, err := authService.UserRepository.FindOneByCondition(ctx, models.User{Email: input.Email})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrAuthRegisterFailed,
			Err:       err,
		}
	}
	if exist != nil {
		return nil, &core.ErrorReturn{
			Err:       nil,
			ErrorCode: response.ErrAuthUserExists,
		}
	}

	password, err := authService.GeneratePassword(input.Password)
	if err != nil {
		return nil, &core.ErrorReturn{
			Err:       err,
			ErrorCode: response.ErrAuthRegisterFailed,
		}
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Status:   enums.StatusActive,
		Password: string(password),
	}
	err = authService.UserRepository.Create(ctx, user)
	if err != nil {
		return nil, &core.ErrorReturn{
			Err:       err,
			ErrorCode: response.ErrAuthRegisterFailed,
		}
	}

	return &types.RegisterOutput{
		User: user,
	}, nil
}

func (authService *AuthService) Logout(ctx context.Context, input types.LogoutInput) *core.ErrorReturn {
	parts := strings.Fields(input.Token)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return &core.ErrorReturn{Err: errors.New("invalid authorization header"), ErrorCode: response.ErrAuthLogoutFailed}
	}
	token := parts[1]
	userID, err := authService.VerifyToken(ctx, token, enums.TokenTypeAccess)
	if err != nil || userID != input.User.ID {
		return &core.ErrorReturn{Err: errors.New("invalid token"), ErrorCode: response.ErrAuthLogoutFailed}
	}
	claims := &models.UserClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.EnvConfig.Auth.JWTSecretKey), nil
	})
	if err != nil || claims.ExpiresAt == nil || claims.ID == "" {
		return &core.ErrorReturn{Err: errors.New("invalid token claims"), ErrorCode: response.ErrAuthLogoutFailed}
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	err = cache.Cache.Set(ctx, fmt.Sprintf(config.CacheKeyBlacklist, userID, claims.ID), true, ttl)
	if err != nil {
		return &core.ErrorReturn{
			Err:       err,
			ErrorCode: response.ErrAuthLogoutFailed,
		}
	}
	return nil
}

func (authService *AuthService) generateToken(ctx context.Context, tokenType enums.TokenType, user *models.User) (string, error) {
	configAuth := config.EnvConfig.Auth
	expiresAt := tokenExpiry(configAuth, tokenType)

	jti := uuid.New().String()
	claims := &models.UserClaims{
		Email:           user.Email,
		Username:        user.Username,
		Type:            tokenType,
		PasswordVersion: passwordVersion(configAuth.JWTSecretKey, user.Password),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   strconv.Itoa(int(user.ID)),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configAuth.JWTSecretKey))
}

func (authService *AuthService) GeneratePassword(plainPassword string) ([]byte, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	return password, nil
}

func (authService *AuthService) VerifyToken(ctx context.Context, token string, tokenType enums.TokenType) (uint, error) {
	id, _, err := authService.VerifyTokenClaims(ctx, token, tokenType)
	return id, err
}

func (authService *AuthService) VerifyTokenClaims(ctx context.Context, token string, tokenType enums.TokenType) (uint, *models.UserClaims, error) {
	configAuth := config.EnvConfig.Auth
	claims := &models.UserClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method")
		}
		return []byte(configAuth.JWTSecretKey), nil
	})
	if err != nil || claims.Type != tokenType || claims.ExpiresAt == nil || claims.ID == "" {
		return 0, nil, errors.New("not valid token parse")
	}
	i, err := strconv.Atoi(claims.Subject)
	if err != nil || i <= 0 {
		return 0, nil, errors.New("invalid token subject")
	}
	userId := uint(i)

	hasBlacklist, err := cache.Cache.Get(ctx, fmt.Sprintf(config.CacheKeyBlacklist, userId, claims.ID))
	if err != nil {
		return 0, nil, err
	}
	if hasBlacklist != nil {
		return 0, nil, errors.New("token invalid")
	}

	return userId, claims, nil
}

func (authService *AuthService) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	user := &models.User{}
	err := authService.UserRepository.FindByID(ctx, id, user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return user, err
}

func (authService *AuthService) GenerateToken(ctx context.Context, tokenType enums.TokenType, user *models.User) (string, error) {
	return authService.generateToken(ctx, tokenType, user)
}

func passwordVersion(secret, passwordHash string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(passwordHash))
	return hex.EncodeToString(mac.Sum(nil))
}

func PasswordVersionForUser(user *models.User) string {
	return passwordVersion(config.EnvConfig.Auth.JWTSecretKey, user.Password)
}

func tokenExpiry(auth config.Auth, tokenType enums.TokenType) time.Time {
	if tokenType == enums.TokenTypeRefresh {
		return time.Now().Add(time.Duration(auth.JWTRefreshExpirationDays) * 24 * time.Hour)
	}
	return time.Now().Add(time.Duration(auth.JWTAccessExpirationMinutes) * time.Minute)
}
