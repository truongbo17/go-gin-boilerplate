package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var CreateRolePermissionsTable = &gormigrate.Migration{
	ID: "2025012201_create_role_permissions_table",
	Migrate: func(tx *gorm.DB) error {
		return tx.Exec(`
			CREATE TABLE role_permissions (
				role_id       INT UNSIGNED NOT NULL,
				permission_id INT UNSIGNED NOT NULL,
				PRIMARY KEY (role_id, permission_id),
				FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
				FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
			);
		`).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Exec("DROP TABLE IF EXISTS role_permissions;").Error
	},
}
