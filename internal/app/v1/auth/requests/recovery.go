package requests

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	"github.com/truongbo17/go-gin-boilerplate/internal/i18n"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (value RegisterRequest) Validate(lang string) error {
	return validation.ValidateStruct(&value,
		validation.Field(&value.Username, required(lang), length(lang, 3, 50)),
		validation.Field(&value.Email,
			required(lang),
			length(lang, 3, 100),
			validation.By(func(input any) error { return plainEmail(lang, input) }),
		),
		validation.Field(&value.Password,
			required(lang),
			length(lang, 8, 64),
			validation.By(func(input any) error { return strongPassword(lang, input) }),
		),
	)
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (value ForgotPasswordRequest) Validate(lang string) error {
	return validation.ValidateStruct(&value,
		validation.Field(&value.Email,
			required(lang),
			length(lang, 3, 100),
			validation.By(func(input any) error { return plainEmail(lang, input) }),
		),
	)
}

type ResetPasswordRequest struct {
	Token           string `json:"token"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (value ResetPasswordRequest) Validate(lang string) error {
	return validation.ValidateStruct(&value,
		validation.Field(&value.Token, required(lang), length(lang, 1, 2048)),
		validation.Field(&value.NewPassword,
			required(lang),
			length(lang, 8, 64),
			validation.By(func(input any) error { return strongPassword(lang, input) }),
		),
		validation.Field(&value.ConfirmPassword, required(lang), validation.By(func(input any) error {
			if input.(string) != value.NewPassword {
				return errors.New(i18n.GetMessage(lang, "validation.password_mismatch", nil))
			}
			return nil
		})),
	)
}

func required(lang string) validation.Rule {
	return validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil))
}

func length(lang string, min, max int) validation.Rule {
	return validation.Length(min, max).Error(i18n.GetMessage(lang, "validation.length", map[string]string{
		"min": fmt.Sprint(min), "max": fmt.Sprint(max),
	}))
}

func plainEmail(lang string, input any) error {
	email := input.(string)
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return errors.New(i18n.GetMessage(lang, "validation.email", nil))
	}
	return nil
}

func strongPassword(lang string, input any) error {
	password := input.(string)
	var upper, lower, digit bool
	for _, r := range password {
		upper = upper || unicode.IsUpper(r)
		lower = lower || unicode.IsLower(r)
		digit = digit || unicode.IsDigit(r)
	}
	for _, check := range []struct {
		ok  bool
		key string
	}{
		{upper, "validation.password_uppercase"},
		{lower, "validation.password_lowercase"},
		{digit, "validation.password_digit"},
		{strings.ContainsAny(password, `!@#$%^&*(),.?":{}|<>`), "validation.password_special"},
	} {
		if !check.ok {
			return errors.New(i18n.GetMessage(lang, check.key, nil))
		}
	}
	return nil
}

func RegisterValidator() gin.HandlerFunc {
	return bindAuthRequest[RegisterRequest]("RegisterRequest")
}

func ForgotPasswordValidator() gin.HandlerFunc {
	return bindAuthRequest[ForgotPasswordRequest]("ForgotPasswordRequest")
}

func ResetPasswordValidator() gin.HandlerFunc {
	return bindAuthRequest[ResetPasswordRequest]("ResetPasswordRequest")
}

func bindAuthRequest[T interface{ Validate(string) error }](key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var value T
		if err := ctx.ShouldBindJSON(&value); err != nil {
			if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
				ctx.AbortWithStatus(http.StatusRequestEntityTooLarge)
				return
			}
			response.ReturnError(ctx, http.StatusUnprocessableEntity, core.ErrUnprocessableEntity,
				i18n.GetMessage(ctx.GetString(config.HeaderLanguage), "error.422", nil), nil)
			return
		}
		if err := value.Validate(ctx.GetString(config.HeaderLanguage)); err != nil {
			response.ReturnError(ctx, http.StatusUnprocessableEntity, core.ErrUnprocessableEntity,
				i18n.GetMessage(ctx.GetString(config.HeaderLanguage), "error.422", nil), []string{err.Error()})
			return
		}
		ctx.Set(key, value)
		ctx.Next()
	}
}
