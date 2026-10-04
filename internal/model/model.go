package models

import (
	"time"

	"gorm.io/gorm"
)

type BasicWithDeleteModel struct {
	BasicModel
	BasicDeleteDateModel
}

type BasicModel struct {
	BasicIDModel
	BasicDateModel
}

type BasicIDModel struct {
	ID uint `json:"id" gorm:"primaryKey,autoIncrement"`
}

type BasicDateModel struct {
	BasicCreateDateModel
	BasicUpdateDateModel
}

type BasicCreateDateModel struct {
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type BasicUpdateDateModel struct {
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type BasicDeleteDateModel struct {
	DeletedAt *time.Time `json:"deleted_at"`
}

type BasicSoftDeleteModel struct {
	DeletedAt *gorm.DeletedAt `json:"deleted_at"`
}
