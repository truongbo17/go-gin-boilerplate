package requests

import (
	"github.com/gin-gonic/gin"
)

type ListUserRequest struct {
	Search  string `form:"search" json:"search"`
	Page    int    `form:"page" json:"page"`
	PerPage int    `form:"per_page" json:"per_page"`
}

func ListUserValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ListUserRequest
		_ = ctx.ShouldBind(&req)

		if req.Page <= 0 {
			req.Page = 1
		}
		if req.PerPage <= 0 {
			req.PerPage = 20
		}

		ctx.Set("ListUserRequest", req)
		ctx.Next()
	}
}
