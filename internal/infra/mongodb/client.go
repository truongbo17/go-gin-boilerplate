package mongodb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Options struct {
	URI         string
	MaxPoolSize uint64
}

// Open connects and checks MongoDB. The caller owns Client.Disconnect(ctx).
func Open(ctx context.Context, opts Options) (*mongo.Client, error) {
	if strings.TrimSpace(opts.URI) == "" {
		return nil, errors.New("MongoDB URI is required")
	}
	clientOptions := options.Client().ApplyURI(opts.URI).SetServerSelectionTimeout(3 * time.Second)
	if opts.MaxPoolSize > 0 {
		clientOptions.SetMaxPoolSize(opts.MaxPoolSize)
	}
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connect MongoDB: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer closeCancel()
		_ = client.Disconnect(closeCtx)
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}
	return client, nil
}
