package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/requests"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/responses"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
)

func (c *AuthController) Register(ctx *gin.Context) {
	requestBody := ctx.MustGet("RegisterRequest").(requests.RegisterRequest)
	result, failure := c.AuthService.Register(ctx.Request.Context(), types.RegisterInput{
		Username: requestBody.Username,
		Email:    requestBody.Email,
		Password: requestBody.Password,
	})
	if failure != nil {
		responses.ReturnError(ctx, failure)
		return
	}
	if result.NotificationErr != nil {
		_ = ctx.Error(result.NotificationErr)
	}
	response.ReturnSuccess(ctx, responses.RegisterResponse{
		ID:          result.User.ID,
		EmailQueued: result.NotificationQueued,
	}, nil)
}

func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	requestBody := ctx.MustGet("ForgotPasswordRequest").(requests.ForgotPasswordRequest)
	if failure := c.AuthService.RequestPasswordReset(ctx.Request.Context(), requestBody.Email); failure != nil {
		responses.ReturnError(ctx, failure)
		return
	}
	responses.ReturnPasswordResetRequested(ctx)
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	requestBody := ctx.MustGet("ResetPasswordRequest").(requests.ResetPasswordRequest)
	failure := c.AuthService.ResetPassword(ctx.Request.Context(), types.ResetPasswordInput{
		Token:    requestBody.Token,
		Password: requestBody.NewPassword,
	})
	if failure != nil {
		responses.ReturnError(ctx, failure)
		return
	}
	response.ReturnSuccess(ctx, nil, nil)
}
