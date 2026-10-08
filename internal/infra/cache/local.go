package cache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"time"
)

type Local struct {
	store *cache.Cache
}

func NewLocal(defaultExpiration, cleanupInterval time.Duration) *Local {
	return &Local{
		store: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (l *Local) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	l.store.Set(key, value, ttl)
	return nil
}

func (l *Local) Get(ctx context.Context, key string) (interface{}, error) {
	val, found := l.store.Get(key)
	if !found {
		return nil, nil
	}
	return val, nil
}

func (l *Local) Delete(ctx context.Context, key string) error {
	l.store.Delete(key)
	return nil
}
