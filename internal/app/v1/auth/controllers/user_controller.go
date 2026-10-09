package controllers

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/requests"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth/responses"
	"github.com/truongbo17/go-gin-boilerplate/internal/request"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type UserController struct {
	Tracer      trace.Tracer
	UserService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		Tracer:      otel.Tracer("UserController"),
		UserService: userService,
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
		responses.ReturnError(ctx, err)
		return
	}

	var userResponses []responses.UserResponse
	for _, user := range users.Items {
		userResponses = append(userResponses, responses.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	response.ReturnSuccess(ctx, userResponses, response.PageMeta(ctx, *users))
}

// GetUserRoles godoc
// @Summary      Get user roles
// @Description  Get all roles assigned to a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  responses.UserRolesResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/users/{id}/roles [get]
// @Security     BearerAuth
func (c *UserController) GetUserRoles(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "GetUserRoles")
	defer span.End()

	id, ok := request.PathID(ctx)
	if !ok {
		return
	}

	roles, err := c.UserService.GetUserRoles(ctxHandler, id)
	if err != nil {
		responses.ReturnError(ctx, err)
		return
	}

	var roleResponses []responses.RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, responses.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Slug:        role.Slug,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	response.ReturnSuccess(ctx, responses.UserRolesResponse{
		Roles: roleResponses,
	}, nil)
}

// AssignRoleToUser godoc
// @Summary      Assign roles to user
// @Description  Assign roles to a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        req  body      requests.AssignRoleToUserRequest  true  "Assign Role Request"
// @Success      200  {object}  response.BaseResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/users/{id}/roles [post]
// @Security     BearerAuth
func (c *UserController) AssignRoleToUser(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "AssignRoleToUser")
	defer span.End()

	id, ok := request.PathID(ctx)
	if !ok {
		return
	}

	var assignRequest, _ = ctx.Get("AssignRoleToUserRequest")
	requestBody, _ := assignRequest.(requests.AssignRoleToUserRequest)

	err := c.UserService.AssignRoleToUser(ctxHandler, types.AssignRoleToUserInput{
		UserID:  id,
		ActorID: ctx.MustGet("user").(*models.User).ID,
		RoleIDs: requestBody.RoleIDs,
	})
	if err != nil {
		responses.ReturnError(ctx, err)
		return
	}

	response.ReturnSuccess(ctx, nil, nil)
}
