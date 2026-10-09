package clickhouse

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Options struct {
	Addresses       []string
	Database        string
	Username        string
	Password        string
	TLS             bool
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Open creates a native-protocol ClickHouse pool and verifies it. The caller
// owns Conn.Close(). Use TLS for managed or public endpoints.
func Open(ctx context.Context, opts Options) (driver.Conn, error) {
	if len(opts.Addresses) == 0 || opts.Database == "" || opts.Username == "" {
		return nil, errors.New("ClickHouse addresses, database, and username are required")
	}
	for _, addr := range opts.Addresses {
		if strings.TrimSpace(addr) == "" {
			return nil, errors.New("ClickHouse address cannot be empty")
		}
	}
	if opts.MaxOpenConns < 0 || opts.MaxIdleConns < 0 || opts.ConnMaxLifetime < 0 || opts.MaxOpenConns > 0 && opts.MaxIdleConns > opts.MaxOpenConns {
		return nil, errors.New("invalid ClickHouse pool settings")
	}
	options := &ch.Options{
		Addr:            opts.Addresses,
		Auth:            ch.Auth{Database: opts.Database, Username: opts.Username, Password: opts.Password},
		DialTimeout:     3 * time.Second,
		MaxOpenConns:    opts.MaxOpenConns,
		MaxIdleConns:    opts.MaxIdleConns,
		ConnMaxLifetime: opts.ConnMaxLifetime,
	}
	if opts.TLS {
		options.TLS = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	conn, err := ch.Open(options)
	if err != nil {
		return nil, fmt.Errorf("connect ClickHouse: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := conn.Ping(pingCtx); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping ClickHouse: %w", err)
	}
	return conn, nil
}
