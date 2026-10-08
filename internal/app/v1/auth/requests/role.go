package requests

import (
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/i18n"
	"github.com/truongbo17/go-gin-boilerplate/internal/request"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

var roleSlugPattern = regexp.MustCompile(`^[a-z][a-z0-9:_-]{0,99}$`)

type ListRoleRequest struct {
	Search string `json:"search" form:"search"`
	request.BaseRequestPaginate
}

func (a ListRoleRequest) Validate(lang string) error {
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

func ListRoleValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var listRoleRequest ListRoleRequest
		if !bindRequest(context, &listRoleRequest) {
			return
		}

		if err := listRoleRequest.Validate(context.GetString(config.HeaderLanguage)); err != nil {
			_ = context.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		context.Set("ListRoleRequest", listRoleRequest)
		context.Next()
	}
}

type CreateRoleRequest struct {
	Name        string `json:"name" form:"name"`
	Slug        string `json:"slug" form:"slug"`
	Description string `json:"description" form:"description"`
}

func (a CreateRoleRequest) Validate(lang string) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.Length(1, 100).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "1", "max": "100"})),
		),
		validation.Field(&a.Slug,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.Length(1, 100).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "1", "max": "100"})),
			validation.Match(roleSlugPattern),
		),
		validation.Field(&a.Description,
			validation.Length(0, 500).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "0", "max": "500"})),
		),
	)
}

func CreateRoleValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var createRoleRequest CreateRoleRequest
		if !bindRequest(context, &createRoleRequest) {
			return
		}

		if err := createRoleRequest.Validate(context.GetString(config.HeaderLanguage)); err != nil {
			_ = context.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		context.Set("CreateRoleRequest", createRoleRequest)
		context.Next()
	}
}

type UpdateRoleRequest struct {
	Name        string `json:"name" form:"name"`
	Slug        string `json:"slug" form:"slug"`
	Description string `json:"description" form:"description"`
}

func (a UpdateRoleRequest) Validate(lang string) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.Length(1, 100).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "1", "max": "100"})),
		),
		validation.Field(&a.Slug,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
			validation.Length(1, 100).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "1", "max": "100"})),
			validation.Match(roleSlugPattern),
		),
		validation.Field(&a.Description,
			validation.Length(0, 500).Error(i18n.GetMessage(lang, "validation.length", map[string]string{"min": "0", "max": "500"})),
		),
	)
}

func UpdateRoleValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var updateRoleRequest UpdateRoleRequest
		if !bindRequest(context, &updateRoleRequest) {
			return
		}

		if err := updateRoleRequest.Validate(context.GetString(config.HeaderLanguage)); err != nil {
			_ = context.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		context.Set("UpdateRoleRequest", updateRoleRequest)
		context.Next()
	}
}

type AssignPermissionToRoleRequest struct {
	PermissionIDs []uint `json:"permission_ids" form:"permission_ids"`
}

func (a AssignPermissionToRoleRequest) Validate(lang string) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.PermissionIDs,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
		),
	)
}

func AssignPermissionToRoleValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var assignPermissionToRoleRequest AssignPermissionToRoleRequest
		if !bindRequest(context, &assignPermissionToRoleRequest) {
			return
		}

		if err := assignPermissionToRoleRequest.Validate(context.GetString(config.HeaderLanguage)); err != nil {
			_ = context.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		context.Set("AssignPermissionToRoleRequest", assignPermissionToRoleRequest)
		context.Next()
	}
}

type AssignRoleToUserRequest struct {
	RoleIDs []uint `json:"role_ids" form:"role_ids"`
}

func (a AssignRoleToUserRequest) Validate(lang string) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.RoleIDs,
			validation.Required.Error(i18n.GetMessage(lang, "validation.required", nil)),
		),
	)
}

func AssignRoleToUserValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var assignRoleToUserRequest AssignRoleToUserRequest
		if !bindRequest(context, &assignRoleToUserRequest) {
			return
		}

		if err := assignRoleToUserRequest.Validate(context.GetString(config.HeaderLanguage)); err != nil {
			_ = context.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		context.Set("AssignRoleToUserRequest", assignRoleToUserRequest)
		context.Next()
	}
}
