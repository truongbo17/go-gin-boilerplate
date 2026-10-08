package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/i18n"
	"github.com/truongbo17/go-gin-boilerplate/internal/request"
	"net/http"
)

type ListPermissionRequest struct {
	Search string `json:"search" form:"search"`
	request.BaseRequestPaginate
}

func (a ListPermissionRequest) Validate(lang string) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.PerPage,
			validation.Min(1).Error(i18n.GetMessage(lang, "validation.min", map[string]string{"min": "1"})),
			validation.Max(100).Error(i18n.GetMessage(lang, "validation.max", map[string]string{"max": "100"})),
		),
		validation.Field(&a.Page,
			validation.Min(0).Error(i18n.GetMessage(lang, "validation.min", map[string]string{"min": "0"})),
			validation.Max(10000),
		),
		validation.Field(&a.Search,
			validation.Length(0, 100).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "0", "max": "100"})),
		),
	)
}

func ListPermissionValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var listPermissionRequest ListPermissionRequest
		if !bindRequest(context, &listPermissionRequest) {
			return
		}

		if err := listPermissionRequest.Validate(context.GetString(config.HeaderLanguage)); err != nil {
			_ = context.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		context.Set("ListPermissionRequest", listPermissionRequest)
		context.Next()
	}
}

type CreatePermissionRequest struct {
	Name        string `json:"name" form:"name"`
	Slug        string `json:"slug" form:"slug"`
	Description string `json:"description" form:"description"`
}

func (a CreatePermissionRequest) Validate(lang string) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.Length(1, 100).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "1", "max": "100"})),
		),
		validation.Field(&a.Slug,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.Length(1, 100).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "1", "max": "100"})),
		),
		validation.Field(&a.Description,
			validation.Length(0, 500).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "0", "max": "500"})),
		),
	)
}

func CreatePermissionValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var createPermissionRequest CreatePermissionRequest
		if !bindRequest(context, &createPermissionRequest) {
			return
		}

		if err := createPermissionRequest.Validate(context.GetString(config.HeaderLanguage)); err != nil {
			_ = context.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		context.Set("CreatePermissionRequest", createPermissionRequest)
		context.Next()
	}
}

type UpdatePermissionRequest struct {
	Name        string `json:"name" form:"name"`
	Slug        string `json:"slug" form:"slug"`
	Description string `json:"description" form:"description"`
}

func (a UpdatePermissionRequest) Validate(lang string) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.Length(1, 100).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "1", "max": "100"})),
		),
		validation.Field(&a.Slug,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.Length(1, 100).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "1", "max": "100"})),
		),
		validation.Field(&a.Description,
			validation.Length(0, 500).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "0", "max": "500"})),
		),
	)
}

func UpdatePermissionValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var updatePermissionRequest UpdatePermissionRequest
		if !bindRequest(context, &updatePermissionRequest) {
			return
		}

		if err := updatePermissionRequest.Validate(context.GetString(config.HeaderLanguage)); err != nil {
			_ = context.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		context.Set("UpdatePermissionRequest", updatePermissionRequest)
		context.Next()
	}
}
