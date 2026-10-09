package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	authcore "github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	coreworker "github.com/truongbo17/go-gin-boilerplate/internal/app/core/worker"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository AuthUserRepository
	TokenBlacklist TokenBlacklist
	AuthConfig     AuthSettings
	Jobs           coreworker.Dispatcher
}

const tokenBlacklistKeyFormat = "token:blacklist:%d_%s"

// AuthSettings contains only the token settings needed by the auth use case.
type AuthSettings struct {
	JWTSecretKey               string
	JWTAccessExpirationMinutes int
	JWTRefreshExpirationDays   int
}

type AuthUserRepository interface {
	FindByID(context.Context, uint) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	Create(context.Context, *models.User) error
	UpdatePassword(context.Context, *models.User, string) error
	UpdatePasswordIfCurrent(context.Context, *models.User, string, string) (bool, error)
}

type TokenBlacklist interface {
	Set(context.Context, string, any, time.Duration) error
	Get(context.Context, string) (any, error)
}

func NewAuthService(userRepository AuthUserRepository, tokenBlacklist TokenBlacklist, authConfig AuthSettings) AuthService {
	return AuthService{UserRepository: userRepository, TokenBlacklist: tokenBlacklist, AuthConfig: authConfig}
}

func (authService *AuthService) ChangePassword(ctx context.Context, user *models.User, oldPassword, newPassword string) (string, *core.ErrorReturn) {
	if user == nil {
		return "", &core.ErrorReturn{ErrorCode: authcore.ErrChangePass}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return "", &core.ErrorReturn{ErrorCode: authcore.ErrChangePass}
	}

	password, err := authService.GeneratePassword(newPassword)
	if err != nil {
		return "", &core.ErrorReturn{ErrorCode: authcore.ErrChangePass, Err: err}
	}

	if err := authService.UserRepository.UpdatePassword(ctx, user, string(password)); err != nil {
		return "", &core.ErrorReturn{ErrorCode: authcore.ErrChangePass, Err: err}
	}

	updatedUser := *user
	updatedUser.Password = string(password)
	token, err := authService.GenerateToken(ctx, enums.TokenTypeAccess, &updatedUser)
	if err != nil {
		return "", &core.ErrorReturn{ErrorCode: authcore.ErrAuthGenerateToken, Err: err}
	}
	user.Password = updatedUser.Password
	return token, nil
}

func (authService *AuthService) Login(ctx context.Context, input types.LoginInput) (*types.LoginOutput, *core.ErrorReturn) {
	user, err := authService.UserRepository.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrAuthLoginFailed,
			Err:       err,
		}
	}
	if user == nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrAuthLoginFailed,
		}
	}
	if user.Status != enums.StatusActive {
		return nil, &core.ErrorReturn{ErrorCode: authcore.ErrAuthLoginFailed}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrAuthLoginFailed,
		}
	}

	token, err := authService.generateToken(ctx, enums.TokenTypeAccess, user)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrAuthGenerateToken,
			Err:       err,
		}
	}

	return &types.LoginOutput{
		AccessToken: token,
	}, nil
}

func (authService *AuthService) Register(ctx context.Context, input types.RegisterInput) (*types.RegisterOutput, *core.ErrorReturn) {
	exist, err := authService.UserRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrAuthRegisterFailed,
			Err:       err,
		}
	}
	if exist != nil {
		return nil, &core.ErrorReturn{
			Err:       nil,
			ErrorCode: authcore.ErrAuthUserExists,
		}
	}
	exist, err = authService.UserRepository.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, &core.ErrorReturn{ErrorCode: authcore.ErrAuthRegisterFailed, Err: err}
	}
	if exist != nil {
		return nil, &core.ErrorReturn{ErrorCode: authcore.ErrAuthUserExists}
	}

	password, err := authService.GeneratePassword(input.Password)
	if err != nil {
		return nil, &core.ErrorReturn{
			Err:       err,
			ErrorCode: authcore.ErrAuthRegisterFailed,
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
		if errors.Is(err, authcore.ErrUserAlreadyExists) {
			return nil, &core.ErrorReturn{ErrorCode: authcore.ErrAuthUserExists}
		}
		return nil, &core.ErrorReturn{
			Err:       err,
			ErrorCode: authcore.ErrAuthRegisterFailed,
		}
	}
	var notificationErr error
	var notificationQueued bool
	if authService.Jobs != nil {
		// A failed notification must not turn a committed registration into a retryable create.
		notificationErr = authService.Jobs.Dispatch(ctx, coreworker.TypeWelcomeEmail, coreworker.WelcomeEmailParams{UserID: user.ID})
		notificationQueued = notificationErr == nil
	}

	return &types.RegisterOutput{
		User:               user,
		NotificationQueued: notificationQueued,
		NotificationErr:    notificationErr,
	}, nil
}

// RequestPasswordReset always has the same public result for known and unknown users.
func (authService *AuthService) RequestPasswordReset(ctx context.Context, email string) *core.ErrorReturn {
	if authService.Jobs == nil {
		return &core.ErrorReturn{ErrorCode: authcore.ErrAuthResetUnavailable}
	}
	if err := authService.Jobs.Dispatch(ctx, coreworker.TypePasswordResetEmail, coreworker.PasswordResetEmailParams{Email: strings.TrimSpace(email)}); err != nil {
		return &core.ErrorReturn{ErrorCode: authcore.ErrAuthResetUnavailable, Err: err}
	}
	return nil
}

func (authService *AuthService) ResetPassword(ctx context.Context, input types.ResetPasswordInput) *core.ErrorReturn {
	id, claims, err := authService.parseTokenClaims(input.Token, enums.TokenTypePasswordReset)
	if err != nil {
		return &core.ErrorReturn{ErrorCode: authcore.ErrAuthResetInvalid}
	}
	user, err := authService.UserRepository.FindByID(ctx, id)
	if err != nil {
		return &core.ErrorReturn{ErrorCode: authcore.ErrAuthResetFailed, Err: err}
	}
	if user == nil || user.Status != enums.StatusActive || claims.PasswordVersion != authService.PasswordVersionForUser(user) {
		return &core.ErrorReturn{ErrorCode: authcore.ErrAuthResetInvalid}
	}
	hash, err := generatePassword(input.Password)
	if err != nil {
		return &core.ErrorReturn{ErrorCode: authcore.ErrAuthResetFailed, Err: err}
	}
	updated, err := authService.UserRepository.UpdatePasswordIfCurrent(ctx, user, user.Password, string(hash))
	if err != nil {
		return &core.ErrorReturn{ErrorCode: authcore.ErrAuthResetFailed, Err: err}
	}
	if !updated {
		return &core.ErrorReturn{ErrorCode: authcore.ErrAuthResetInvalid}
	}
	return nil
}

func (authService *AuthService) Logout(ctx context.Context, input types.LogoutInput) *core.ErrorReturn {
	parts := strings.Fields(input.Token)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return &core.ErrorReturn{Err: errors.New("invalid authorization header"), ErrorCode: authcore.ErrAuthLogoutFailed}
	}
	token := parts[1]
	userID, err := authService.VerifyToken(ctx, token, enums.TokenTypeAccess)
	if err != nil || userID != input.User.ID {
		return &core.ErrorReturn{Err: errors.New("invalid token"), ErrorCode: authcore.ErrAuthLogoutFailed}
	}
	claims := &types.UserClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method")
		}
		return []byte(authService.AuthConfig.JWTSecretKey), nil
	})
	if err != nil || claims.ExpiresAt == nil || claims.ID == "" {
		return &core.ErrorReturn{Err: errors.New("invalid token claims"), ErrorCode: authcore.ErrAuthLogoutFailed}
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	err = authService.TokenBlacklist.Set(ctx, fmt.Sprintf(tokenBlacklistKeyFormat, userID, claims.ID), true, ttl)
	if err != nil {
		return &core.ErrorReturn{
			Err:       err,
			ErrorCode: authcore.ErrAuthLogoutFailed,
		}
	}
	return nil
}

func (authService *AuthService) generateToken(ctx context.Context, tokenType enums.TokenType, user *models.User) (string, error) {
	configAuth := authService.AuthConfig
	expiresAt := tokenExpiry(configAuth, tokenType)

	jti := uuid.New().String()
	claims := &types.UserClaims{
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
	return generatePassword(plainPassword)
}

func generatePassword(plainPassword string) ([]byte, error) {
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

func (authService *AuthService) VerifyTokenClaims(ctx context.Context, token string, tokenType enums.TokenType) (uint, *types.UserClaims, error) {
	userId, claims, err := authService.parseTokenClaims(token, tokenType)
	if err != nil {
		return 0, nil, err
	}

	hasBlacklist, err := authService.TokenBlacklist.Get(ctx, fmt.Sprintf(tokenBlacklistKeyFormat, userId, claims.ID))
	if err != nil {
		return 0, nil, err
	}
	if hasBlacklist != nil {
		return 0, nil, errors.New("token invalid")
	}

	return userId, claims, nil
}

func (authService *AuthService) parseTokenClaims(token string, tokenType enums.TokenType) (uint, *types.UserClaims, error) {
	configAuth := authService.AuthConfig
	claims := &types.UserClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
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

	return userId, claims, nil
}

func (authService *AuthService) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	return authService.UserRepository.FindByID(ctx, id)
}

func (authService *AuthService) GenerateToken(ctx context.Context, tokenType enums.TokenType, user *models.User) (string, error) {
	return authService.generateToken(ctx, tokenType, user)
}

func passwordVersion(secret, passwordHash string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(passwordHash))
	return hex.EncodeToString(mac.Sum(nil))
}

func (authService *AuthService) PasswordVersionForUser(user *models.User) string {
	return passwordVersion(authService.AuthConfig.JWTSecretKey, user.Password)
}

func tokenExpiry(auth AuthSettings, tokenType enums.TokenType) time.Time {
	if tokenType == enums.TokenTypePasswordReset {
		return time.Now().Add(15 * time.Minute)
	}
	if tokenType == enums.TokenTypeRefresh {
		return time.Now().Add(time.Duration(auth.JWTRefreshExpirationDays) * 24 * time.Hour)
	}
	return time.Now().Add(time.Duration(auth.JWTAccessExpirationMinutes) * time.Minute)
}
