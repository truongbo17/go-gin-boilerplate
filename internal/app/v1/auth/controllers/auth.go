package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/requests"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/responses"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthController struct {
	Tracer      trace.Tracer
	AuthService services.AuthService
	UserService services.UserService
}

func NewAuthController() *AuthController {
	return &AuthController{
		Tracer:      otel.Tracer("AuthController"),
		AuthService: services.NewAuthService(),
		UserService: services.NewUserService(),
	}
}

// Login godoc
// @Summary      Login
// @Description  login a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      requests.LoginRequest true "Login Request"
// @Success      200  {object}  responses.LoginResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "Login")
	defer span.End()

	var loginRequest, _ = ctx.Get("LoginRequest")
	requestBody, _ := loginRequest.(requests.LoginRequest)

	token, err := c.AuthService.Login(ctxHandler, types.LoginInput{
		Username: requestBody.Username,
		Password: requestBody.Password,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	response.ReturnSuccess(ctx, responses.LoginResponse{
		AccessToken: token.AccessToken,
	}, nil)
	return
}

// Logout godoc
// @Summary      Logout
// @Description  logout a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.BaseResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/auth/logout [post]
// @Security BearerAuth
func (c *AuthController) Logout(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "Logout")
	defer span.End()

	user := ctx.MustGet("user").(*models.User)
	err := c.AuthService.Logout(ctxHandler, types.LogoutInput{
		User:  user,
		Token: ctx.GetHeader(config.HeaderAuth),
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	response.ReturnSuccess(ctx, nil, nil)
}

// Me godoc
// @Summary      Me
// @Description  Get me
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.MeResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/auth/me [get]
// @Security BearerAuth
func (c *AuthController) Me(ctx *gin.Context) {
	_, span := c.Tracer.Start(ctx.Request.Context(), "Me")
	defer span.End()

	user := ctx.MustGet("user").(*models.User)
	response.ReturnSuccess(ctx, responses.MeResponse{
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		CreateAt:  user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil)
	return
}

func (c *AuthController) ChangePass(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "ChangePass")
	defer span.End()

	var changePasswordRequest, _ = ctx.Get("ChangePasswordRequest")
	requestBody, _ := changePasswordRequest.(requests.ChangePasswordRequest)
	user := ctx.MustGet("user").(*models.User)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.OldPassword))
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, response.ErrChangePass, nil)
		return
	}
	newPass, err := c.AuthService.GeneratePassword(requestBody.NewPassword)
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, response.ErrChangePass, err)
		return
	}

	err = c.UserService.UserRepository.Update(ctxHandler, user, &models.User{
		Password: string(newPass),
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, response.ErrChangePass, err)
		return
	}
	user.Password = string(newPass)
	accessToken, err := c.AuthService.GenerateToken(ctxHandler, enums.TokenTypeAccess, user)
	if err != nil {
		response.ReturnError(ctx, http.StatusInternalServerError, response.ErrAuthGenerateToken, err)
		return
	}

	response.ReturnSuccess(ctx, responses.LoginResponse{
		AccessToken: accessToken,
	}, nil)
}
