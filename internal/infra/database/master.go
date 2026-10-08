package database

import (
	"context"
	"fmt"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/logger"
	"net"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectMaster() error {
	configDB := config.EnvConfig.Master
	if configDB.Username != "" {
		dsn := (&mysqldriver.Config{
			User: configDB.Username, Passwd: configDB.Password,
			Net: "tcp", Addr: net.JoinHostPort(configDB.Host, configDB.Port),
			DBName: configDB.Database, ParseTime: true, Loc: time.Local,
			InterpolateParams: true, Params: map[string]string{"charset": "utf8mb4"},
			Timeout: 3 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second,
		}).FormatDSN()

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
					ParameterizedQueries:      true,
					Colorful:                  false,
				},
			)
		}

		db, err := gorm.Open(gormmysql.Open(dsn), configMysql)
		if err != nil {
			return fmt.Errorf("connect MySQL: %w", err)
		}

		dbConfig, err := db.DB()
		if err != nil {
			return fmt.Errorf("get MySQL connection pool: %w", err)
		}
		dbConfig.SetMaxOpenConns(30)
		dbConfig.SetMaxIdleConns(15)
		dbConfig.SetConnMaxLifetime(15 * time.Minute)
		dbConfig.SetConnMaxIdleTime(5 * time.Minute)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := dbConfig.PingContext(ctx); err != nil {
			_ = dbConfig.Close()
			return fmt.Errorf("ping MySQL: %w", err)
		}

		if err = db.Use(otelgorm.NewPlugin()); err != nil {
			_ = dbConfig.Close()
			return fmt.Errorf("instrument MySQL: %w", err)
		}

		DB = db

		fmt.Println("Success connected to Mysql")
	}
	return nil
}
