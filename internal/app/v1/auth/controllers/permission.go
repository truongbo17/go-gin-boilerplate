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

type PermissionController struct {
	Tracer            trace.Tracer
	PermissionService services.PermissionService
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		Tracer:            otel.Tracer("PermissionController"),
		PermissionService: services.NewPermissionService(),
	}
}

// ListPermissions godoc
// @Summary      List permissions
// @Description  Get paginated list of permissions
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        req  query     requests.ListPermissionRequest  true  "List Permissions Request"
// @Success      200  {object}  responses.PermissionResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/permissions [get]
// @Security     BearerAuth
func (c *PermissionController) ListPermissions(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "ListPermissions")
	defer span.End()

	var listRequest, _ = ctx.Get("ListPermissionRequest")
	requestBody, _ := listRequest.(requests.ListPermissionRequest)

	permissions, err := c.PermissionService.ListPermissions(ctxHandler, types.ListPermissionsInput{
		Search:  requestBody.Search,
		Page:    requestBody.Page,
		PerPage: requestBody.PerPage,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	var permissionResponses []responses.PermissionResponse
	for _, permission := range *permissions.Data {
		permissionResponses = append(permissionResponses, responses.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Slug:        permission.Slug,
			Description: permission.Description,
			CreatedAt:   permission.CreatedAt,
			UpdatedAt:   permission.UpdatedAt,
		})
	}

	response.ReturnSuccess(ctx, permissionResponses, permissions.MetaData)
}

// CreatePermission godoc
// @Summary      Create a new permission
// @Description  Create a new permission
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        req  body      requests.CreatePermissionRequest  true  "Create Permission Request"
// @Success      200  {object}  responses.PermissionResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/permissions [post]
// @Security     BearerAuth
func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "CreatePermission")
	defer span.End()

	var createRequest, _ = ctx.Get("CreatePermissionRequest")
	requestBody, _ := createRequest.(requests.CreatePermissionRequest)

	permission, err := c.PermissionService.CreatePermission(ctxHandler, types.CreatePermissionInput{
		Name:        requestBody.Name,
		Slug:        requestBody.Slug,
		Description: requestBody.Description,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	response.ReturnSuccess(ctx, responses.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Slug:        permission.Slug,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}, nil)
}

// UpdatePermission godoc
// @Summary      Update a permission
// @Description  Update an existing permission
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Permission ID"
// @Param        req  body      requests.UpdatePermissionRequest  true  "Update Permission Request"
// @Success      200  {object}  responses.PermissionResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/permissions/{id} [put]
// @Security     BearerAuth
func (c *PermissionController) UpdatePermission(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "UpdatePermission")
	defer span.End()

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var updateRequest, _ = ctx.Get("UpdatePermissionRequest")
	requestBody, _ := updateRequest.(requests.UpdatePermissionRequest)

	permission, err := c.PermissionService.UpdatePermission(ctxHandler, types.UpdatePermissionInput{
		ID:          id,
		Name:        requestBody.Name,
		Slug:        requestBody.Slug,
		Description: requestBody.Description,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	response.ReturnSuccess(ctx, responses.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Slug:        permission.Slug,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}, nil)
}

// DeletePermission godoc
// @Summary      Delete a permission
// @Description  Delete an existing permission
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Permission ID"
// @Success      200  {object}  response.BaseResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /api/v1/rbac/permissions/{id} [delete]
// @Security     BearerAuth
func (c *PermissionController) DeletePermission(ctx *gin.Context) {
	ctxHandler, span := c.Tracer.Start(ctx.Request.Context(), "DeletePermission")
	defer span.End()

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	err := c.PermissionService.DeletePermission(ctxHandler, types.DeletePermissionInput{
		ID: id,
	})
	if err != nil {
		response.ReturnError(ctx, http.StatusOK, err.ErrorCode, err.Err)
		return
	}

	response.ReturnSuccess(ctx, nil, nil)
}
