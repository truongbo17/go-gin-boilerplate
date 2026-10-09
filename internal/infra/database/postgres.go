package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresOptions configures an optional PostgreSQL connection. DSN can be a
// URL or a libpq connection string; keep credentials out of logs.
type PostgresOptions struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// OpenPostgres opens and verifies a PostgreSQL pool. The caller owns DB().Close().
func OpenPostgres(ctx context.Context, opts PostgresOptions) (*gorm.DB, error) {
	if strings.TrimSpace(opts.DSN) == "" {
		return nil, errors.New("PostgreSQL DSN is required")
	}
	if opts.MaxOpenConns < 0 || opts.MaxIdleConns < 0 || opts.MaxIdleConns > 0 && opts.MaxOpenConns > 0 && opts.MaxIdleConns > opts.MaxOpenConns {
		return nil, errors.New("invalid PostgreSQL pool size")
	}
	if opts.ConnMaxLifetime < 0 || opts.ConnMaxIdleTime < 0 {
		return nil, errors.New("PostgreSQL pool durations cannot be negative")
	}

	db, err := gorm.Open(postgres.Open(opts.DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
	})
	if err != nil {
		return nil, fmt.Errorf("connect PostgreSQL: %w", err)
	}
	pool, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get PostgreSQL pool: %w", err)
	}
	if opts.MaxOpenConns > 0 {
		pool.SetMaxOpenConns(opts.MaxOpenConns)
	}
	if opts.MaxIdleConns > 0 {
		pool.SetMaxIdleConns(opts.MaxIdleConns)
	}
	pool.SetConnMaxLifetime(opts.ConnMaxLifetime)
	pool.SetConnMaxIdleTime(opts.ConnMaxIdleTime)
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.PingContext(pingCtx); err != nil {
		_ = pool.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		_ = pool.Close()
		return nil, fmt.Errorf("instrument PostgreSQL: %w", err)
	}
	return db, nil
}
