package controllers

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/requests"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/responses"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type UserController struct {
	Tracer      trace.Tracer
	UserService services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		Tracer:      otel.Tracer("UserController"),
		UserService: services.NewUserService(),
	}
}

// ListUsers godoc
// @Summary      List users
// @Description  Get paginated list of users with search
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        req  query     requests.ListUserRequest  true  "List Users Request"
// @Success      200  {object}  responses.UserResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/users [get]
// @Security     BearerAuth
func (c *UserController) ListUsers(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "ListUsers")
	defer span.End()

	var listRequest, _ = ctx.Get("ListUserRequest")
	requestBody, _ := listRequest.(requests.ListUserRequest)

	users, err := c.UserService.ListUsers(ctxHandler, types.ListUsersInput{
		Search:  requestBody.Search,
		Page:    requestBody.Page,
		PerPage: requestBody.PerPage,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	var userResponses []responses.UserResponse
	for _, user := range *users.Data {
		userResponses = append(userResponses, responses.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	response.ReturnSuccess(ctx, userResponses, users.MetaData)
}
