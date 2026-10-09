package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"net/url"

	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
)

var ErrInvalidMailAddress = errors.New("invalid mail address")

type MailUsers interface {
	FindByID(context.Context, uint) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
}

type MailSender interface {
	Send(context.Context, string, string, string) error
}

type ResetTokenGenerator interface {
	GenerateToken(context.Context, enums.TokenType, *models.User) (string, error)
}

type MailService struct {
	Users    MailUsers
	Tokens   ResetTokenGenerator
	Sender   MailSender
	ResetURL string
}

func (service MailService) SendWelcome(ctx context.Context, userID uint) error {
	user, err := service.Users.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("find welcome email recipient: %w", err)
	}
	if user == nil {
		return nil
	}
	if !validMailAddress(user.Email) {
		return ErrInvalidMailAddress
	}
	return service.Sender.Send(ctx, user.Email, "Welcome", "Welcome, "+user.Username+"!\n")
}

func (service MailService) SendPasswordReset(ctx context.Context, email string) error {
	if !validMailAddress(email) {
		return ErrInvalidMailAddress
	}
	user, err := service.Users.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("find password reset recipient: %w", err)
	}
	if user == nil || user.Status != enums.StatusActive {
		return nil
	}
	if !validMailAddress(user.Email) {
		return ErrInvalidMailAddress
	}
	token, err := service.Tokens.GenerateToken(ctx, enums.TokenTypePasswordReset, user)
	if err != nil {
		return fmt.Errorf("generate password reset token: %w", err)
	}
	link, err := url.Parse(service.ResetURL)
	if err != nil || link.Host == "" || (link.Scheme != "http" && link.Scheme != "https") {
		return errors.New("invalid password reset URL")
	}
	query := link.Query()
	query.Set("token", token)
	link.RawQuery = query.Encode()
	return service.Sender.Send(ctx, user.Email, "Reset your password", "Open this link within 15 minutes to reset your password:\n"+link.String()+"\n")
}

func validMailAddress(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}
