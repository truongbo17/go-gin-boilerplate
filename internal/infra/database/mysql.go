package database

import (
	"context"
	"fmt"
	"net"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func Open(ctx context.Context, configDB config.Master, env string, log *logrus.Logger) (*gorm.DB, error) {
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
		if (env == config.DebugMode || env == config.LocalMode) && log != nil {
			configMysql.Logger = gormLogger.New(
				log,
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
			return nil, fmt.Errorf("connect MySQL: %w", err)
		}

		dbConfig, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("get MySQL connection pool: %w", err)
		}
		dbConfig.SetMaxOpenConns(configDB.MaxOpenConns)
		dbConfig.SetMaxIdleConns(configDB.MaxIdleConns)
		dbConfig.SetConnMaxLifetime(time.Duration(configDB.ConnMaxLifetimeMinutes) * time.Minute)
		dbConfig.SetConnMaxIdleTime(time.Duration(configDB.ConnMaxIdleTimeMinutes) * time.Minute)
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := dbConfig.PingContext(ctx); err != nil {
			_ = dbConfig.Close()
			return nil, fmt.Errorf("ping MySQL: %w", err)
		}

		if err = db.Use(otelgorm.NewPlugin()); err != nil {
			_ = dbConfig.Close()
			return nil, fmt.Errorf("instrument MySQL: %w", err)
		}

		return db, nil
	}
	return nil, fmt.Errorf("MySQL configuration is missing a username")
}
