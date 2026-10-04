package cache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type Local struct {
	store *cache.Cache
	locks sync.Map
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

func (l *Local) Increment(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	lock := l.getLock(key)
	lock.Lock()
	defer lock.Unlock()

	val, found := l.store.Get(key)

	var timeValue int64
	if found {
		timeValue = val.(int64) + 1
	} else {
		timeValue = 1
	}

	err := l.Set(ctx, key, timeValue, ttl)
	if err != nil {
		return 0, err
	}

	return timeValue, nil
}

func (l *Local) getLock(key string) *sync.Mutex {
	actual, _ := l.locks.LoadOrStore(key, &sync.Mutex{})
	return actual.(*sync.Mutex)
}
