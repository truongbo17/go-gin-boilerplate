package services

import (
	"context"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	authcore "github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/page"
	"strings"
)

type RoleService struct {
	RoleRepository RoleStore
}

type RoleStore interface {
	ListPage(context.Context, page.Query) (*page.Result[models.Role], error)
	Create(context.Context, *models.Role) error
	FindByID(context.Context, uint) (*models.Role, error)
	UpdateFields(context.Context, *models.Role, types.UpdateRoleInput) error
	Delete(context.Context, uint) error
}

func NewRoleService(roleRepository RoleStore) RoleService {
	return RoleService{RoleRepository: roleRepository}
}

func (rs *RoleService) ListRoles(ctx context.Context, input types.ListRolesInput) (*page.Result[models.Role], *core.ErrorReturn) {
	roles, err := rs.RoleRepository.ListPage(ctx, page.Query{Number: input.Page, Size: input.PerPage, Search: input.Search})
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrRoleInternalError,
			Err:       err,
		}
	}

	return roles, nil
}

func (rs *RoleService) CreateRole(ctx context.Context, input types.CreateRoleInput) (*models.Role, *core.ErrorReturn) {
	if isReservedAdminSlug(input.Slug) {
		return nil, &core.ErrorReturn{ErrorCode: core.ErrForbidden}
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
			ErrorCode: authcore.ErrRoleCreateFailed,
			Err:       err,
		}
	}

	return role, nil
}

func (rs *RoleService) UpdateRole(ctx context.Context, input types.UpdateRoleInput) (*models.Role, *core.ErrorReturn) {
	if isReservedAdminSlug(input.Slug) {
		return nil, &core.ErrorReturn{ErrorCode: core.ErrForbidden}
	}
	role, err := rs.RoleRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrRoleInternalError,
			Err:       err,
		}
	}
	if role == nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrRoleNotFound,
		}
	}
	if isReservedAdminSlug(role.Slug) {
		return nil, &core.ErrorReturn{ErrorCode: core.ErrForbidden}
	}

	err = rs.RoleRepository.UpdateFields(ctx, role, input)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrRoleUpdateFailed,
			Err:       err,
		}
	}

	return role, nil
}

func (rs *RoleService) DeleteRole(ctx context.Context, input types.DeleteRoleInput) *core.ErrorReturn {
	role, err := rs.RoleRepository.FindByID(ctx, input.ID)
	if err != nil {
		return &core.ErrorReturn{ErrorCode: authcore.ErrRoleInternalError, Err: err}
	}
	if role == nil {
		return &core.ErrorReturn{ErrorCode: authcore.ErrRoleNotFound}
	}
	if isReservedAdminSlug(role.Slug) {
		return &core.ErrorReturn{ErrorCode: core.ErrForbidden}
	}
	err = rs.RoleRepository.Delete(ctx, input.ID)
	if err != nil {
		return &core.ErrorReturn{
			ErrorCode: authcore.ErrRoleDeleteFailed,
			Err:       err,
		}
	}

	return nil
}

func isReservedAdminSlug(slug string) bool {
	return strings.EqualFold(strings.TrimSpace(slug), "admin")
}
