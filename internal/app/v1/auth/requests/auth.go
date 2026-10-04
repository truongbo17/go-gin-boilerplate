package requests

import (
	"errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/i18n"
	"net/http"
	"regexp"
)

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (a LoginRequest) Validate(lang string) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Username, validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil))),
		validation.Field(&a.Password,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.Length(8, 64).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "8", "max": "64"})),
		),
	)
}

func LoginValidator() gin.HandlerFunc {
	return func(context *gin.Context) {

		var loginRequest LoginRequest
		_ = context.ShouldBind(&loginRequest)

		if err := loginRequest.Validate(context.GetString(config.HeaderLanguage)); err != nil {
			_ = context.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		context.Set("LoginRequest", loginRequest)

		context.Next()
	}
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" form:"old_password"`
	NewPassword     string `json:"new_password" form:"new_password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}

func (a ChangePasswordRequest) Validate(lang string) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.OldPassword,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
		),
		validation.Field(&a.NewPassword,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.Length(8, 64).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "8", "max": "64"})),
			validation.Match(regexp.MustCompile(`[A-Z]`)).Error(i18n.GetMessage(lang, "validation.password_uppercase", nil)),
			validation.Match(regexp.MustCompile(`[a-z]`)).Error(i18n.GetMessage(lang, "validation.password_lowercase", nil)),
			validation.Match(regexp.MustCompile(`\d`)).Error(i18n.GetMessage(lang, "validation.password_digit", nil)),
			validation.Match(regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)).Error(i18n.GetMessage(lang, "validation.password_special", nil)),
		),
		validation.Field(&a.ConfirmPassword,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.By(func(value interface{}) error {
				if value.(string) != a.NewPassword {
					return errors.New(i18n.GetMessage(lang, "validation.password_mismatch", nil))
				}
				return nil
			}),
		),
	)
}

func ChangePasswordValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var changePasswordRequest ChangePasswordRequest
		_ = context.ShouldBind(&changePasswordRequest)

		if err := changePasswordRequest.Validate(context.GetString(config.HeaderLanguage)); err != nil {
			_ = context.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		context.Set("ChangePasswordRequest", changePasswordRequest)

		context.Next()
	}
}
