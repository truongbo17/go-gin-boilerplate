package services

import (
	"context"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	authcore "github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/page"
)

type PermissionService struct {
	PermissionRepository PermissionStore
}

type PermissionStore interface {
	ListPage(context.Context, page.Query) (*page.Result[models.Permission], error)
	Create(context.Context, *models.Permission) error
	FindByID(context.Context, uint) (*models.Permission, error)
	UpdateFields(context.Context, *models.Permission, types.UpdatePermissionInput) error
	Delete(context.Context, uint) error
	CheckPermission(context.Context, uint, string) (bool, error)
	GetRolePermissions(context.Context, uint) ([]models.Permission, error)
	AssignPermissionToRole(context.Context, uint, []uint) error
}

func NewPermissionService(permissionRepository PermissionStore) PermissionService {
	return PermissionService{PermissionRepository: permissionRepository}
}

func (ps *PermissionService) CheckPermission(ctx context.Context, userID uint, permissionSlug string) (bool, error) {
	return ps.PermissionRepository.CheckPermission(ctx, userID, permissionSlug)
}

func (ps *PermissionService) ListPermissions(ctx context.Context, input types.ListPermissionsInput) (*page.Result[models.Permission], *core.ErrorReturn) {
	permissions, err := ps.PermissionRepository.ListPage(ctx, page.Query{Number: input.Page, Size: input.PerPage, Search: input.Search})
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrPermissionInternalError,
			Err:       err,
		}
	}

	return permissions, nil
}

func (ps *PermissionService) CreatePermission(ctx context.Context, input types.CreatePermissionInput) (*models.Permission, *core.ErrorReturn) {
	permission := &models.Permission{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
	}

	err := ps.PermissionRepository.Create(ctx, permission)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrPermissionCreateFailed,
			Err:       err,
		}
	}

	return permission, nil
}

func (ps *PermissionService) UpdatePermission(ctx context.Context, input types.UpdatePermissionInput) (*models.Permission, *core.ErrorReturn) {
	permission, err := ps.PermissionRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrPermissionInternalError,
			Err:       err,
		}
	}
	if permission == nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrPermissionNotFound,
		}
	}

	err = ps.PermissionRepository.UpdateFields(ctx, permission, input)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrPermissionUpdateFailed,
			Err:       err,
		}
	}

	return permission, nil
}

func (ps *PermissionService) DeletePermission(ctx context.Context, input types.DeletePermissionInput) *core.ErrorReturn {
	err := ps.PermissionRepository.Delete(ctx, input.ID)
	if err != nil {
		return &core.ErrorReturn{
			ErrorCode: authcore.ErrPermissionDeleteFailed,
			Err:       err,
		}
	}

	return nil
}

func (ps *PermissionService) GetRolePermissions(ctx context.Context, roleID uint) ([]models.Permission, *core.ErrorReturn) {
	permissions, err := ps.PermissionRepository.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrPermissionInternalError,
			Err:       err,
		}
	}
	if permissions == nil {
		permissions = []models.Permission{}
	}
	return permissions, nil
}

func (ps *PermissionService) AssignPermissionToRole(ctx context.Context, input types.AssignPermissionToRoleInput) *core.ErrorReturn {
	if err := ps.PermissionRepository.AssignPermissionToRole(ctx, input.RoleID, input.PermissionIDs); err != nil {
		return &core.ErrorReturn{ErrorCode: authcore.ErrRoleAssignPermission, Err: err}
	}
	return nil
}
