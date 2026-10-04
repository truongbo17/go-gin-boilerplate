package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/logger"
)

func Migrate() {
	migrations := []*gormigrate.Migration{CreateUsersTable, CreateRolesTable, CreatePermissionsTable, CreateRolePermissionsTable, CreateUserRolesTable}
	m := gormigrate.New(database.DB, &gormigrate.Options{TableName: config.TableMigrate, IDColumnName: "id", IDColumnSize: config.DefaultStringSizeMySql, UseTransaction: true, ValidateUnknownMigrations: true}, migrations)
	if err := m.Migrate(); err != nil {
		panic(err)
	}
	logger.LogrusLogger.Infoln("Migration successful")
}
