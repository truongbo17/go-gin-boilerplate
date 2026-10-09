package config

import (
	"errors"
	"net/mail"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Mail struct {
	Enabled  bool   `mapstructure:"MAIL_ENABLED"`
	Host     string `mapstructure:"MAIL_HOST"`
	Port     string `mapstructure:"MAIL_PORT"`
	Username string `mapstructure:"MAIL_USERNAME"`
	Password string `mapstructure:"MAIL_PASSWORD"`
	From     string `mapstructure:"MAIL_FROM"`
	TLSMode  string `mapstructure:"MAIL_TLS_MODE"`
	ResetURL string `mapstructure:"PASSWORD_RESET_URL"`
}

func (config *Config) validateMail() error {
	if !config.Mail.Enabled {
		return nil
	}
	if err := validation.ValidateStruct(config,
		validation.Field(&config.Cache.RedisHost, validation.Required, is.Host),
		validation.Field(&config.Cache.RedisPort, validation.Required, is.Port),
		validation.Field(&config.Mail.Host, validation.Required),
		validation.Field(&config.Mail.Port, validation.Required, is.Port),
		validation.Field(&config.Mail.From, validation.Required),
		validation.Field(&config.Mail.TLSMode, validation.In("none", "starttls", "implicit")),
		validation.Field(&config.Mail.ResetURL, validation.Required),
	); err != nil {
		return err
	}
	if config.App.Env == ReleaseMode && config.Mail.TLSMode == "none" {
		return errors.New("MAIL_TLS_MODE=none is only allowed outside release mode")
	}
	address, err := mail.ParseAddress(config.Mail.From)
	if err != nil || address.Address != config.Mail.From {
		return errors.New("MAIL_FROM must be a plain email address")
	}
	resetURL, err := url.Parse(config.Mail.ResetURL)
	if err != nil || resetURL.Host == "" || (resetURL.Scheme != "https" && !(config.App.Env != ReleaseMode && resetURL.Scheme == "http")) {
		return errors.New("PASSWORD_RESET_URL must be an absolute HTTPS URL (HTTP allowed outside release mode)")
	}
	return nil
}
