package models

type UserRole struct {
	UserID uint `gorm:"column:user_id;type:INT UNSIGNED;not null;" json:"user_id"`
	RoleID uint `gorm:"column:role_id;type:INT UNSIGNED;not null;" json:"role_id"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
