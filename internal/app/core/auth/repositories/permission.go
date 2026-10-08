package repositories

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	repository.BaseRepository[models.Permission]
}

func NewPermissionRepository() PermissionRepository {
	return PermissionRepository{BaseRepository: repository.NewBaseRepository[models.Permission](database.DB)}
}

func (r *PermissionRepository) WithTransaction(tx *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		BaseRepository: *r.BaseRepository.WithTransaction(tx),
	}
}
