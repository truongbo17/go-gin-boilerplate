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

type RoleController struct {
	Tracer            trace.Tracer
	RoleService       services.RoleService
	PermissionService services.PermissionService
}

func NewRoleController() *RoleController {
	return &RoleController{
		Tracer:            otel.Tracer("RoleController"),
		RoleService:       services.NewRoleService(),
		PermissionService: services.NewPermissionService(),
	}
}

// ListRoles godoc
// @Summary      List roles
// @Description  Get paginated list of roles
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        req  query     requests.ListRoleRequest  true  "List Roles Request"
// @Success      200  {object}  responses.RoleResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/roles [get]
// @Security     BearerAuth
func (c *RoleController) ListRoles(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "ListRoles")
	defer span.End()

	var listRequest, _ = ctx.Get("ListRoleRequest")
	requestBody, _ := listRequest.(requests.ListRoleRequest)

	roles, err := c.RoleService.ListRoles(ctxHandler, types.ListRolesInput{
		Search:  requestBody.Search,
		Page:    requestBody.Page,
		PerPage: requestBody.PerPage,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	var roleResponses []responses.RoleResponse
	for _, role := range *roles.Data {
		roleResponses = append(roleResponses, responses.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Slug:        role.Slug,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	response.ReturnSuccess(ctx, roleResponses, roles.MetaData)
}

// CreateRole godoc
// @Summary      Create a new role
// @Description  Create a new role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        req  body      requests.CreateRoleRequest  true  "Create Role Request"
// @Success      200  {object}  responses.RoleResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/roles [post]
// @Security     BearerAuth
func (c *RoleController) CreateRole(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "CreateRole")
	defer span.End()

	var createRequest, _ = ctx.Get("CreateRoleRequest")
	requestBody, _ := createRequest.(requests.CreateRoleRequest)

	role, err := c.RoleService.CreateRole(ctxHandler, types.CreateRoleInput{
		Name:        requestBody.Name,
		Slug:        requestBody.Slug,
		Description: requestBody.Description,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	response.ReturnSuccess(ctx, responses.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Slug:        role.Slug,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil)
}

// UpdateRole godoc
// @Summary      Update a role
// @Description  Update an existing role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Role ID"
// @Param        req  body      requests.UpdateRoleRequest  true  "Update Role Request"
// @Success      200  {object}  responses.RoleResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/roles/{id} [put]
// @Security     BearerAuth
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "UpdateRole")
	defer span.End()

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var updateRequest, _ = ctx.Get("UpdateRoleRequest")
	requestBody, _ := updateRequest.(requests.UpdateRoleRequest)

	role, err := c.RoleService.UpdateRole(ctxHandler, types.UpdateRoleInput{
		ID:          id,
		Name:        requestBody.Name,
		Slug:        requestBody.Slug,
		Description: requestBody.Description,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	response.ReturnSuccess(ctx, responses.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Slug:        role.Slug,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil)
}

// DeleteRole godoc
// @Summary      Delete a role
// @Description  Delete an existing role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  response.BaseResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/roles/{id} [delete]
// @Security     BearerAuth
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "DeleteRole")
	defer span.End()

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	err := c.RoleService.DeleteRole(ctxHandler, types.DeleteRoleInput{
		ID: id,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	response.ReturnSuccess(ctx, nil, nil)
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
func (c *RoleController) GetUserRoles(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "GetUserRoles")
	defer span.End()

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	roles, err := c.PermissionService.GetUserRoles(ctxHandler, id)
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
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

// GetRolePermissions godoc
// @Summary      Get role permissions
// @Description  Get all permissions assigned to a role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  responses.RolePermissionsResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/roles/{id}/permissions [get]
// @Security     BearerAuth
func (c *RoleController) GetRolePermissions(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "GetRolePermissions")
	defer span.End()

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	permissions, err := c.PermissionService.GetRolePermissions(ctxHandler, id)
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	var permissionResponses []responses.PermissionResponse
	for _, permission := range permissions {
		permissionResponses = append(permissionResponses, responses.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Slug:        permission.Slug,
			Description: permission.Description,
			CreatedAt:   permission.CreatedAt,
			UpdatedAt:   permission.UpdatedAt,
		})
	}

	response.ReturnSuccess(ctx, responses.RolePermissionsResponse{
		Permissions: permissionResponses,
	}, nil)
}

// AssignPermissionToRole godoc
// @Summary      Assign permissions to role
// @Description  Assign permissions to a role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Role ID"
// @Param        req  body      requests.AssignPermissionToRoleRequest  true  "Assign Permission Request"
// @Success      200  {object}  response.BaseResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/roles/{id}/permissions [post]
// @Security     BearerAuth
func (c *RoleController) AssignPermissionToRole(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "AssignPermissionToRole")
	defer span.End()

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var assignRequest, _ = ctx.Get("AssignPermissionToRoleRequest")
	requestBody, _ := assignRequest.(requests.AssignPermissionToRoleRequest)

	err := c.PermissionService.AssignPermissionToRole(ctxHandler, types.AssignPermissionToRoleInput{
		RoleID:        id,
		PermissionIDs: requestBody.PermissionIDs,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	response.ReturnSuccess(ctx, nil, nil)
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
func (c *RoleController) AssignRoleToUser(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "AssignRoleToUser")
	defer span.End()

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var assignRequest, _ = ctx.Get("AssignRoleToUserRequest")
	requestBody, _ := assignRequest.(requests.AssignRoleToUserRequest)

	err := c.PermissionService.AssignRoleToUser(ctxHandler, types.AssignRoleToUserInput{
		UserID:  id,
		RoleIDs: requestBody.RoleIDs,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	response.ReturnSuccess(ctx, nil, nil)
}
