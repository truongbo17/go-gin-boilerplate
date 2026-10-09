package requests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/i18n"
	"github.com/truongbo17/go-gin-boilerplate/internal/request"
)

type ListUserRequest struct {
	Search  string `form:"search" json:"search"`
	Page    int    `form:"page" json:"page"`
	PerPage int    `form:"per_page" json:"per_page"`
}

func ListUserValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ListUserRequest
		if !request.Bind(ctx, &req) {
			return
		}

		if req.Page <= 0 {
			req.Page = 1
		}
		if req.Page > 10000 {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		if req.PerPage <= 0 {
			req.PerPage = 20
		}
		if req.PerPage > 100 {
			req.PerPage = 100
		}
		if len(req.Search) > 100 {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		ctx.Set("ListUserRequest", req)
		ctx.Next()
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
		if !request.Bind(context, &assignRoleToUserRequest) {
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
