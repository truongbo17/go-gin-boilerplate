package database

import (
	"fmt"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/logger"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectMaster() {
	configDB := config.EnvConfig.Master
	if configDB.Username != "" {
		address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true",
			configDB.Username,
			configDB.Password,
			configDB.Host,
			configDB.Port,
			configDB.Database,
		)

		configMysql := &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 nil,
		}
		if config.EnvConfig.Env == config.DebugMode || config.EnvConfig.Env == config.LocalMode {
			//configMysql.Logger = gormLogger.Default.LogMode(gormLogger.Info)
			configMysql.Logger = gormLogger.New(
				logger.LogrusLogger,
				gormLogger.Config{
					SlowThreshold:             time.Second,
					LogLevel:                  gormLogger.Info,
					IgnoreRecordNotFoundError: true,
					Colorful:                  false,
				},
			)
		}

		db, err := gorm.Open(mysql.Open(address), configMysql)
		if err != nil {
			panic("Connected failed, check your MySql")
		}

		dbConfig, _ := db.DB()
		dbConfig.SetMaxOpenConns(30)
		dbConfig.SetMaxIdleConns(15)
		dbConfig.SetConnMaxLifetime(15 * time.Minute)
		dbConfig.SetConnMaxIdleTime(5 * time.Minute)

		if err = db.Use(otelgorm.NewPlugin()); err != nil {
			panic(err)
		}

		DB = db

		fmt.Println("Success connected to Mysql")
	}
}
