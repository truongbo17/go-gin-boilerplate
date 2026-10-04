package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var CreateUsersTable = &gormigrate.Migration{
	ID: "2025012201_create_user_table",
	Migrate: func(tx *gorm.DB) error {
		return tx.Exec(`
			CREATE TABLE users (
				id                 INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
				username           VARCHAR(50) UNIQUE          NOT NULL,
				email              VARCHAR(100) UNIQUE         NOT NULL,
				password           CHAR(64)                    NOT NULL,
				status             TINYINT DEFAULT 1,
				created_by         INT UNSIGNED DEFAULT 0,
				updated_by         INT UNSIGNED DEFAULT 0,
				created_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at         DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
				deleted_at         DATETIME NULL
			);
		`).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Exec("DROP TABLE IF EXISTS users;").Error
	},
}
