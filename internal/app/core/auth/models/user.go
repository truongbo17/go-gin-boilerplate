package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	models "github.com/truongbo17/go-gin-boilerplate/internal/model"
)

type User struct {
	models.BasicModel
	models.BasicSoftDeleteModel

	Username  string       `gorm:"column:username;type:VARCHAR(50);not null;" json:"username"`
	Email     string       `gorm:"column:email;type:VARCHAR(100);not null;" json:"email"`
	Password  string       `gorm:"column:password;type:CHAR(64);not null;" json:"-"`
	Status    enums.Status `gorm:"column:status;type:TINYINT;default:1;" json:"status"`
	CreatedBy uint         `gorm:"column:created_by;type:INT UNSIGNED;default:0;" json:"created_by"`
	UpdatedBy uint         `gorm:"column:updated_by;type:INT UNSIGNED;default:0;" json:"updated_by"`
}

func (User) TableName() string {
	return "users"
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username        string          `json:"username"`
	Email           string          `json:"email"`
	Type            enums.TokenType `json:"type"`
	PasswordVersion string          `json:"password_version"`
	Permissions     []string        `json:"permissions,omitempty"`
}
