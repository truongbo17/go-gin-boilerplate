package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var CreateRolesTable = &gormigrate.Migration{
	ID: "2025012201_create_roles_table",
	Migrate: func(tx *gorm.DB) error {
		return tx.Exec(`
			CREATE TABLE roles (
				id                 INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
				name               VARCHAR(100) NOT NULL,
				slug               VARCHAR(100) UNIQUE NOT NULL,
				description        TEXT,
				created_by         INT UNSIGNED DEFAULT 0,
				updated_by         INT UNSIGNED DEFAULT 0,
				created_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at         DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
				deleted_at         DATETIME NULL
			);
		`).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Exec("DROP TABLE IF EXISTS roles;").Error
	},
}
