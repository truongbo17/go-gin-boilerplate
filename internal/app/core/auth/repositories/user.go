package repositories

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"gorm.io/gorm"
	"sync"
)

type UserRepository struct {
	repository.BaseRepository[models.User]
}

var (
	userRepo         UserRepository
	userRepoRepoOnce sync.Once
)

func NewUserRepository() UserRepository {
	userRepoRepoOnce.Do(func() {
		db := database.DB
		userRepo = UserRepository{
			BaseRepository: repository.NewBaseRepository[models.User](db),
		}
	})
	return userRepo
}

func (r *UserRepository) WithTransaction(tx *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: *r.BaseRepository.WithTransaction(tx),
	}
}
