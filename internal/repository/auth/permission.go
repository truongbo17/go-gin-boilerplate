package auth

import (
	"context"

	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/page"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"gorm.io/gorm"
)

func (r *PermissionRepository) ListPage(ctx context.Context, query page.Query) (*page.Result[models.Permission], error) {
	params := repository.NewPaginateParam()
	params.Page, params.PerPage = query.Number, query.Size
	if query.Search != "" {
		params.ExtraWheres = append(params.ExtraWheres, repository.WhereClause{
			Query: "name LIKE ? OR slug LIKE ?",
			Args:  []any{"%" + query.Search + "%", "%" + query.Search + "%"},
		})
	}
	return r.Paginate(ctx, params)
}

type PermissionRepository struct {
	repository.BaseRepository[models.Permission]
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{BaseRepository: repository.NewBaseRepository[models.Permission](db)}
}

func (r *PermissionRepository) FindByID(ctx context.Context, id uint) (*models.Permission, error) {
	return r.FindOneByCondition(ctx, map[string]any{"id": id})
}

func (r *PermissionRepository) UpdateFields(ctx context.Context, permission *models.Permission, input types.UpdatePermissionInput) error {
	if err := r.Update(ctx, permission, map[string]any{
		"name": input.Name, "slug": input.Slug, "description": input.Description,
	}); err != nil {
		return err
	}
	permission.Name, permission.Slug, permission.Description = input.Name, input.Slug, input.Description
	return nil
}

func (r *PermissionRepository) WithTransaction(tx *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		BaseRepository: *r.BaseRepository.WithTransaction(tx),
	}
}

func (r *PermissionRepository) CheckPermission(ctx context.Context, userID uint, permissionSlug string) (bool, error) {
	db := r.DB.WithContext(ctx)
	admin, err := hasAdminRole(db, userID)
	if err != nil || admin {
		return admin, err
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

func (r *PermissionRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.DB.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) AssignPermissionToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}
		for _, permissionID := range permissionIDs {
			if err := tx.Create(&models.RolePermission{RoleID: roleID, PermissionID: permissionID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
