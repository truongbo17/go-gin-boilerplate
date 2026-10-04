package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var CreateUserRolesTable = &gormigrate.Migration{
	ID: "2025012201_create_user_roles_table",
	Migrate: func(tx *gorm.DB) error {
		return tx.Exec(`
			CREATE TABLE user_roles (
				user_id INT UNSIGNED NOT NULL,
				role_id INT UNSIGNED NOT NULL,
				PRIMARY KEY (user_id, role_id),
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
				FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
			);
		`).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Exec("DROP TABLE IF EXISTS user_roles;").Error
	},
}
