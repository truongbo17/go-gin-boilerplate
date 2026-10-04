package repositories

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"gorm.io/gorm"
	"sync"
)

type RoleRepository struct {
	repository.BaseRepository[models.Role]
}

var (
	roleRepo         RoleRepository
	roleRepoRepoOnce sync.Once
)

func NewRoleRepository() RoleRepository {
	roleRepoRepoOnce.Do(func() {
		db := database.DB
		roleRepo = RoleRepository{
			BaseRepository: repository.NewBaseRepository[models.Role](db),
		}
	})
	return roleRepo
}

func (r *RoleRepository) WithTransaction(tx *gorm.DB) *RoleRepository {
	return &RoleRepository{
		BaseRepository: *r.BaseRepository.WithTransaction(tx),
	}
}
