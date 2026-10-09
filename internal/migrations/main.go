package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	migrations := []*gormigrate.Migration{CreateUsersTable, CreateRolesTable, CreatePermissionsTable, CreateRolePermissionsTable, CreateUserRolesTable}
	m := gormigrate.New(db, &gormigrate.Options{TableName: config.TableMigrate, IDColumnName: "id", IDColumnSize: config.DefaultStringSizeMySql, UseTransaction: true, ValidateUnknownMigrations: true}, migrations)
	return m.Migrate()
}
