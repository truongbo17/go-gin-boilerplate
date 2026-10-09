package auth

import (
	"context"
	"errors"
	"fmt"

	mysqldriver "github.com/go-sql-driver/mysql"
	authcore "github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/page"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"gorm.io/gorm"
)

func (r *UserRepository) ListPage(ctx context.Context, query page.Query) (*page.Result[models.User], error) {
	params := repository.NewPaginateParam()
	params.Page, params.PerPage = query.Number, query.Size
	if query.Search != "" {
		params.ExtraWheres = append(params.ExtraWheres, repository.WhereClause{
			Query: "username LIKE ? OR email LIKE ?",
			Args:  []any{"%" + query.Search + "%", "%" + query.Search + "%"},
		})
	}
	return r.Paginate(ctx, params)
}

type UserRepository struct {
	repository.BaseRepository[models.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{BaseRepository: repository.NewBaseRepository[models.User](db)}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return normalizeUserCreateError(r.BaseRepository.Create(ctx, user))
}

func normalizeUserCreateError(err error) error {
	if err == nil {
		return nil
	}
	var mysqlErr *mysqldriver.MySQLError
	if errors.Is(err, gorm.ErrDuplicatedKey) || (errors.As(err, &mysqlErr) && mysqlErr.Number == 1062) {
		return fmt.Errorf("%w: %v", authcore.ErrUserAlreadyExists, err)
	}
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	return r.FindOneByCondition(ctx, map[string]any{"id": id})
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	return r.FindOneByCondition(ctx, map[string]any{"username": username})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.FindOneByCondition(ctx, map[string]any{"email": email})
}

func (r *UserRepository) UpdatePassword(ctx context.Context, user *models.User, passwordHash string) error {
	return r.Update(ctx, user, map[string]any{"password": passwordHash})
}

// Compare-and-swap prevents two concurrent reset requests from consuming one token.
func (r *UserRepository) UpdatePasswordIfCurrent(ctx context.Context, user *models.User, currentHash, nextHash string) (bool, error) {
	result := r.DB.WithContext(ctx).Model(&models.User{}).
		Where("id = ? AND password = ?", user.ID, currentHash).
		Update("password", nextHash)
	return result.RowsAffected == 1, result.Error
}

func (r *UserRepository) WithTransaction(tx *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: *r.BaseRepository.WithTransaction(tx),
	}
}

func (r *UserRepository) CreateWithAdminRole(ctx context.Context, user *models.User, admin bool) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		if !admin {
			return nil
		}

		role := models.Role{Slug: "admin", Name: "Administrator"}
		if err := tx.Where("slug = ?", role.Slug).FirstOrCreate(&role).Error; err != nil {
			return err
		}
		assignment := models.UserRole{UserID: user.ID, RoleID: role.ID}
		return tx.Where("user_id = ? AND role_id = ?", assignment.UserID, assignment.RoleID).
			FirstOrCreate(&assignment).Error
	})
}

func (r *UserRepository) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.DB.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *UserRepository) AssignRoleToUser(ctx context.Context, actorID, userID uint, roleIDs []uint) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		actorIsAdmin, err := hasAdminRole(tx, actorID)
		if err != nil {
			return err
		}
		if !actorIsAdmin {
			targetIsAdmin, err := hasAdminRole(tx, userID)
			if err != nil {
				return err
			}
			var selectedAdminRoles int64
			if len(roleIDs) > 0 {
				if err := tx.Model(&models.Role{}).Where("id IN ? AND slug = ?", roleIDs, "admin").Count(&selectedAdminRoles).Error; err != nil {
					return err
				}
			}
			if targetIsAdmin || selectedAdminRoles > 0 {
				return gorm.ErrInvalidData
			}
		}
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			if err := tx.Create(&models.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
