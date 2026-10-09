package auth

import (
	"context"

	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/page"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"gorm.io/gorm"
)

func (r *RoleRepository) ListPage(ctx context.Context, query page.Query) (*page.Result[models.Role], error) {
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

type RoleRepository struct {
	repository.BaseRepository[models.Role]
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{BaseRepository: repository.NewBaseRepository[models.Role](db)}
}

func (r *RoleRepository) FindByID(ctx context.Context, id uint) (*models.Role, error) {
	return r.FindOneByCondition(ctx, map[string]any{"id": id})
}

func (r *RoleRepository) UpdateFields(ctx context.Context, role *models.Role, input types.UpdateRoleInput) error {
	if err := r.Update(ctx, role, map[string]any{
		"name": input.Name, "slug": input.Slug, "description": input.Description, "updated_by": input.UserID,
	}); err != nil {
		return err
	}
	role.Name, role.Slug, role.Description, role.UpdatedBy = input.Name, input.Slug, input.Description, input.UserID
	return nil
}

func (r *RoleRepository) WithTransaction(tx *gorm.DB) *RoleRepository {
	return &RoleRepository{
		BaseRepository: *r.BaseRepository.WithTransaction(tx),
	}
}
