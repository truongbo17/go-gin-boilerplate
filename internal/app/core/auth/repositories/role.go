package repositories

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"gorm.io/gorm"
)

type RoleRepository struct {
	repository.BaseRepository[models.Role]
}

func NewRoleRepository() RoleRepository {
	return RoleRepository{BaseRepository: repository.NewBaseRepository[models.Role](database.DB)}
}

func (r *RoleRepository) WithTransaction(tx *gorm.DB) *RoleRepository {
	return &RoleRepository{
		BaseRepository: *r.BaseRepository.WithTransaction(tx),
	}
}
