package services

import (
	"context"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/repositories"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"

	"gorm.io/gorm"
)

type PermissionService struct {
	PermissionRepository repositories.PermissionRepository
}

func NewPermissionService() PermissionService {
	return PermissionService{PermissionRepository: repositories.NewPermissionRepository()}
}

func (ps *PermissionService) CheckPermission(ctx context.Context, userID uint, permissionSlug string) (bool, error) {
	db := database.DB.WithContext(ctx)
	admin, err := hasAdminRole(db, userID)
	if err != nil {
		return false, err
	}
	if admin {
		return true, nil
	}
	var count int64
	err = db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("user_roles.user_id = ? AND permissions.slug = ? AND roles.deleted_at IS NULL AND permissions.deleted_at IS NULL", userID, permissionSlug).
		Count(&count).Error
	return count > 0, err
}

func (ps *PermissionService) ListPermissions(ctx context.Context, input types.ListPermissionsInput) (*response.PaginateResponse[models.Permission], *core.ErrorReturn) {
	params := repository.NewPaginateParam()
	params.Page = input.Page
	params.PerPage = input.PerPage
	params.Path = "/api/v1/rbac/permissions"

	if input.Search != "" {
		params.ExtraWheres = append(params.ExtraWheres, repository.WhereClause{
			Query: "name LIKE ? OR slug LIKE ?",
			Args:  []interface{}{"%" + input.Search + "%", "%" + input.Search + "%"},
		})
	}

	permissions, err := ps.PermissionRepository.Paginate(ctx, params)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrPermissionInternalError,
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
			ErrorCode: response.ErrPermissionCreateFailed,
			Err:       err,
		}
	}

	return permission, nil
}

func (ps *PermissionService) UpdatePermission(ctx context.Context, input types.UpdatePermissionInput) (*models.Permission, *core.ErrorReturn) {
	permission, err := ps.PermissionRepository.FindOneByCondition(ctx, map[string]interface{}{
		"id": input.ID,
	})
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrPermissionInternalError,
			Err:       err,
		}
	}
	if permission == nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrPermissionNotFound,
		}
	}

	updateData := map[string]interface{}{
		"name":        input.Name,
		"slug":        input.Slug,
		"description": input.Description,
	}

	err = ps.PermissionRepository.Update(ctx, permission, updateData)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrPermissionUpdateFailed,
			Err:       err,
		}
	}

	return permission, nil
}

func (ps *PermissionService) DeletePermission(ctx context.Context, input types.DeletePermissionInput) *core.ErrorReturn {
	err := ps.PermissionRepository.Delete(ctx, input.ID)
	if err != nil {
		return &core.ErrorReturn{
			ErrorCode: response.ErrPermissionDeleteFailed,
			Err:       err,
		}
	}

	return nil
}

func (ps *PermissionService) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, *core.ErrorReturn) {
	var userRoles []models.UserRole
	if err := database.DB.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrRoleInternalError,
			Err:       err,
		}
	}

	if len(userRoles) == 0 {
		return []models.Role{}, nil
	}

	roleIDs := make([]uint, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	var roles []models.Role
	if err := database.DB.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrRoleInternalError,
			Err:       err,
		}
	}

	return roles, nil
}

func (ps *PermissionService) GetRolePermissions(ctx context.Context, roleID uint) ([]models.Permission, *core.ErrorReturn) {
	var rolePermissions []models.RolePermission
	if err := database.DB.Where("role_id = ?", roleID).Find(&rolePermissions).Error; err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrPermissionInternalError,
			Err:       err,
		}
	}

	if len(rolePermissions) == 0 {
		return []models.Permission{}, nil
	}

	permissionIDs := make([]uint, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissionIDs[i] = rp.PermissionID
	}

	var permissions []models.Permission
	if err := database.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrPermissionInternalError,
			Err:       err,
		}
	}

	return permissions, nil
}

func (ps *PermissionService) AssignPermissionToRole(ctx context.Context, input types.AssignPermissionToRoleInput) *core.ErrorReturn {
	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", input.RoleID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}
		for _, permissionID := range input.PermissionIDs {
			if err := tx.Create(&models.RolePermission{RoleID: input.RoleID, PermissionID: permissionID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return &core.ErrorReturn{ErrorCode: response.ErrRoleAssignPermission, Err: err}
	}
	return nil
}

func (ps *PermissionService) AssignRoleToUser(ctx context.Context, input types.AssignRoleToUserInput) *core.ErrorReturn {
	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		actorIsAdmin, err := hasAdminRole(tx, input.ActorID)
		if err != nil {
			return err
		}
		if !actorIsAdmin {
			targetIsAdmin, err := hasAdminRole(tx, input.UserID)
			if err != nil {
				return err
			}
			var selectedAdminRoles int64
			if len(input.RoleIDs) > 0 {
				if err := tx.Model(&models.Role{}).Where("id IN ? AND slug = ?", input.RoleIDs, "admin").Count(&selectedAdminRoles).Error; err != nil {
					return err
				}
			}
			if targetIsAdmin || selectedAdminRoles > 0 {
				return gorm.ErrInvalidData
			}
		}
		if err := tx.Where("user_id = ?", input.UserID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}
		for _, roleID := range input.RoleIDs {
			if err := tx.Create(&models.UserRole{UserID: input.UserID, RoleID: roleID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return &core.ErrorReturn{ErrorCode: response.ErrRoleAssignUser, Err: err}
	}
	return nil
}

func hasAdminRole(db *gorm.DB, userID uint) (bool, error) {
	var count int64
	err := db.Model(&models.Role{}).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.slug = ?", userID, "admin").
		Count(&count).Error
	return count > 0, err
}
