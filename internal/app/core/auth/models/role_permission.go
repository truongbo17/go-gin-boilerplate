package models

type RolePermission struct {
	RoleID       uint `gorm:"column:role_id;type:INT UNSIGNED;not null;" json:"role_id"`
	PermissionID uint `gorm:"column:permission_id;type:INT UNSIGNED;not null;" json:"permission_id"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
