package services

import (
	"context"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/repositories"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
	"strings"
	"sync"
)

type RoleService struct {
	RoleRepository repositories.RoleRepository
}

var (
	roleService     RoleService
	roleServiceOnce sync.Once
)

func NewRoleService() RoleService {
	roleServiceOnce.Do(func() {
		roleService = RoleService{
			RoleRepository: repositories.NewRoleRepository(),
		}
	})

	return roleService
}

func (rs *RoleService) ListRoles(ctx context.Context, input types.ListRolesInput) (*response.PaginateResponse[models.Role], *core.ErrorReturn) {
	params := repository.NewPaginateParam()
	params.Page = input.Page
	params.PerPage = input.PerPage
	params.Path = "/api/v1/rbac/roles"

	if input.Search != "" {
		params.ExtraWheres = append(params.ExtraWheres, repository.WhereClause{
			Query: "name LIKE ? OR slug LIKE ?",
			Args:  []interface{}{"%" + input.Search + "%", "%" + input.Search + "%"},
		})
	}

	roles, err := rs.RoleRepository.Paginate(ctx, params)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrRoleInternalError,
			Err:       err,
		}
	}

	return roles, nil
}

func (rs *RoleService) CreateRole(ctx context.Context, input types.CreateRoleInput) (*models.Role, *core.ErrorReturn) {
	if isReservedAdminSlug(input.Slug) {
		return nil, &core.ErrorReturn{ErrorCode: response.ErrForbidden}
	}
	role := &models.Role{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		CreatedBy:   input.UserID,
	}

	err := rs.RoleRepository.Create(ctx, role)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrRoleCreateFailed,
			Err:       err,
		}
	}

	return role, nil
}

func (rs *RoleService) UpdateRole(ctx context.Context, input types.UpdateRoleInput) (*models.Role, *core.ErrorReturn) {
	if isReservedAdminSlug(input.Slug) {
		return nil, &core.ErrorReturn{ErrorCode: response.ErrForbidden}
	}
	role, err := rs.RoleRepository.FindOneByCondition(ctx, map[string]interface{}{
		"id": input.ID,
	})
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrRoleInternalError,
			Err:       err,
		}
	}
	if role == nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrRoleNotFound,
		}
	}
	if isReservedAdminSlug(role.Slug) {
		return nil, &core.ErrorReturn{ErrorCode: response.ErrForbidden}
	}

	updateData := map[string]interface{}{
		"name":        input.Name,
		"slug":        input.Slug,
		"description": input.Description,
		"updated_by":  input.UserID, // Set updated_by to the user updating the role
	}

	err = rs.RoleRepository.Update(ctx, role, updateData)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrRoleUpdateFailed,
			Err:       err,
		}
	}

	return role, nil
}

func (rs *RoleService) DeleteRole(ctx context.Context, input types.DeleteRoleInput) *core.ErrorReturn {
	role, err := rs.RoleRepository.FindOneByCondition(ctx, map[string]interface{}{"id": input.ID})
	if err != nil {
		return &core.ErrorReturn{ErrorCode: response.ErrRoleInternalError, Err: err}
	}
	if role == nil {
		return &core.ErrorReturn{ErrorCode: response.ErrRoleNotFound}
	}
	if isReservedAdminSlug(role.Slug) {
		return &core.ErrorReturn{ErrorCode: response.ErrForbidden}
	}
	err = rs.RoleRepository.Delete(ctx, input.ID)
	if err != nil {
		return &core.ErrorReturn{
			ErrorCode: response.ErrRoleDeleteFailed,
			Err:       err,
		}
	}

	return nil
}

func isReservedAdminSlug(slug string) bool {
	return strings.EqualFold(strings.TrimSpace(slug), "admin")
}
