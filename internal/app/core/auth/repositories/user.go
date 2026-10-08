package repositories

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	repository.BaseRepository[models.User]
}

func NewUserRepository() UserRepository {
	return UserRepository{BaseRepository: repository.NewBaseRepository[models.User](database.DB)}
}

func (r *UserRepository) WithTransaction(tx *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: *r.BaseRepository.WithTransaction(tx),
	}
}
