package models

import (
	models "github.com/truongbo17/go-gin-boilerplate/internal/model"
)

type Role struct {
	models.BasicModel
	models.BasicSoftDeleteModel

	Name        string `gorm:"column:name;type:VARCHAR(100);not null;" json:"name"`
	Slug        string `gorm:"column:slug;type:VARCHAR(100);unique;not null;" json:"slug"`
	Description string `gorm:"column:description;type:TEXT;" json:"description"`
	CreatedBy   uint   `gorm:"column:created_by;type:INT UNSIGNED;default:0;" json:"created_by"`
	UpdatedBy   uint   `gorm:"column:updated_by;type:INT UNSIGNED;default:0;" json:"updated_by"`
}

func (Role) TableName() string {
	return "roles"
}
