package auth

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"gorm.io/gorm"
)

func hasAdminRole(db *gorm.DB, userID uint) (bool, error) {
	var count int64
	err := db.Model(&models.Role{}).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.slug = ?", userID, "admin").
		Count(&count).Error
	return count > 0, err
}
