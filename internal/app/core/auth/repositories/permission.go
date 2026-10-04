package repositories

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"gorm.io/gorm"
	"sync"
)

type PermissionRepository struct {
	repository.BaseRepository[models.Permission]
}

var (
	permissionRepo         PermissionRepository
	permissionRepoRepoOnce sync.Once
)

func NewPermissionRepository() PermissionRepository {
	permissionRepoRepoOnce.Do(func() {
		db := database.DB
		permissionRepo = PermissionRepository{
			BaseRepository: repository.NewBaseRepository[models.Permission](db),
		}
	})
	return permissionRepo
}

func (r *PermissionRepository) WithTransaction(tx *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		BaseRepository: *r.BaseRepository.WithTransaction(tx),
	}
}
