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
	DB           *gorm.DB
	Redis        *redisclient.Client
	Mail         config.Mail
	JWTSecretKey string
}

// NewHandler wires the auth worker adapter to its core mail service.
func NewHandler(params Dependencies) (*Handler, error) {
	if params.DB == nil || params.Redis == nil {
		return nil, errors.New("auth mail worker requires a database and Redis")
	}
	users := authrepository.NewUserRepository(params.DB)
	return &Handler{Mail: services.MailService{
		Users: users,
		Tokens: new(services.NewAuthService(users, cache.NewRedis(params.Redis), services.AuthSettings{
			JWTSecretKey: params.JWTSecretKey,
		})),
		Sender:   mail.Sender{Config: params.Mail},
		ResetURL: params.Mail.ResetURL,
	}}, nil
}
