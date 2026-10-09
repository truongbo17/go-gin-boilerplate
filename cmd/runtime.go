package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	redisclient "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/cache"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/logger"
	redisinfra "github.com/truongbo17/go-gin-boilerplate/internal/infra/redis"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/tracer"
	"gorm.io/gorm"
)

type requirements struct {
	database bool
	redis    bool
	cache    bool
	tracing  bool
	worker   bool
}

// runtime owns process resources for one command invocation.
type runtime struct {
	config config.Config
	log    *logrus.Logger
	db     *gorm.DB
	redis  *redisclient.Client
	cache  cache.ICache
	close  []func() error
}

func openRuntime(ctx context.Context, needs requirements) (result *runtime, err error) {
	var loaded *config.Config
	if needs.worker {
		loaded, err = config.LoadWorker()
	} else if needs.database && !needs.cache {
		loaded, err = config.LoadDatabase()
	} else {
		loaded, err = config.Load()
	}
	if err != nil {
		return nil, err
	}
	if needs.worker && loaded.Mail.Enabled {
		needs.database = true
	}
	r := &runtime{config: *loaded}
	defer func() {
		if err != nil {
			err = errors.Join(err, r.Close())
		}
	}()

	log, logFile := logger.Open(loaded.App.Env)
	r.log = log
	r.close = append(r.close, logFile.Close)
	if needs.tracing {
		provider, startErr := tracer.Start(loaded.Tracer)
		if startErr != nil {
			return nil, startErr
		}
		if provider != nil {
			r.close = append(r.close, func() error {
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				return provider.Shutdown(shutdownCtx)
			})
		}
	}
	if needs.database {
		r.db, err = database.Open(ctx, loaded.Master, loaded.App.Env, r.log)
		if err != nil {
			return nil, err
		}
		pool, poolErr := r.db.DB()
		if poolErr != nil {
			return nil, fmt.Errorf("get MySQL pool: %w", poolErr)
		}
		r.close = append(r.close, pool.Close)
	}
	if needs.redis || loaded.Mail.Enabled || (needs.cache && loaded.Cache.CacheStore == config.CacheStoreRedis) {
		r.redis, err = redisinfra.Open(loaded.Cache)
		if err != nil {
			return nil, err
		}
		if r.redis == nil {
			return nil, errors.New("Redis is required for this command")
		}
		r.close = append(r.close, r.redis.Close)
	}
	if needs.cache {
		r.cache = cache.NewStore(loaded.Cache, r.redis)
	}
	return r, nil
}

func (r *runtime) Close() error {
	var closeErrors []error
	for i := len(r.close) - 1; i >= 0; i-- {
		if err := r.close[i](); err != nil {
			closeErrors = append(closeErrors, err)
		}
	}
	r.close = nil
	return errors.Join(closeErrors...)
}

func withRuntime(cmd *cobra.Command, needs requirements, run func(*runtime) error) (err error) {
	r, err := openRuntime(cmd.Context(), needs)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, r.Close()) }()
	return run(r)
}

func withDatabase(cmd *cobra.Command, run func(*gorm.DB) error) error {
	return withRuntime(cmd, requirements{database: true}, func(r *runtime) error {
		return run(r.db)
	})
}
