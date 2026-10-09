package auth

import (
	"errors"

	redisclient "github.com/redis/go-redis/v9"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/cache"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/mail"
	authrepository "github.com/truongbo17/go-gin-boilerplate/internal/repository/auth"
	"gorm.io/gorm"
)

type Dependencies struct {
	Mail         config.Mail
	JWTSecretKey string
	DB           *gorm.DB
	Redis        *redisclient.Client
}

// NewHandler wires the auth worker adapter to its core mail service.
func NewHandler(deps Dependencies) (*Handler, error) {
	if deps.DB == nil || deps.Redis == nil {
		return nil, errors.New("auth mail worker requires a database and Redis")
	}
	users := authrepository.NewUserRepository(deps.DB)
	return &Handler{Mail: services.MailService{
		Users: users,
		Tokens: new(services.NewAuthService(users, cache.NewRedis(deps.Redis), services.AuthSettings{
			JWTSecretKey: deps.JWTSecretKey,
		})),
		Sender:   mail.Sender{Config: deps.Mail},
		ResetURL: deps.Mail.ResetURL,
	}}, nil
}
