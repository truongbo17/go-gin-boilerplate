package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/cmd/cli"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/cache"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	httpclient "github.com/truongbo17/go-gin-boilerplate/internal/infra/http"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/logger"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/redis"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/tracer"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/worker/client"
)

var rootCmd = &cobra.Command{
	Use:   "ggb",
	Short: "Reusable Gin API base",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "version" {
			return nil
		}
		if err := config.Init(); err != nil {
			return err
		}
		if config.EnvConfig.App.Env == config.ReleaseMode {
			gin.SetMode(gin.ReleaseMode)
		}
		logger.Init()
		if err := tracer.InitTracerOTEL(); err != nil {
			return err
		}
		if err := redis.ConnectRedis(); err != nil {
			return err
		}
		cache.InitCache()
		if err := database.ConnectDatabase(); err != nil {
			return err
		}
		httpclient.InitBaseRequest()
		if config.EnvConfig.Cache.CacheStore == config.CacheStoreRedis {
			client.InitClient()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(StartServerCmd, StartWorkerCmd, cli.VersionCmd, cli.MigrateCmd, cli.CreateUserCmd)
}

func Execute() (err error) {
	defer func() {
		var cleanupErrors []error
		if httpclient.Request != nil {
			httpclient.Request.Client.CloseIdleConnections()
		}
		if client.WorkerClient != nil {
			if closeErr := client.WorkerClient.Close(); closeErr != nil {
				cleanupErrors = append(cleanupErrors, fmt.Errorf("close worker client: %w", closeErr))
			}
		}
		if redis.ClientRedis != nil {
			if closeErr := redis.ClientRedis.Close(); closeErr != nil {
				cleanupErrors = append(cleanupErrors, fmt.Errorf("close Redis: %w", closeErr))
			}
		}
		if database.DB != nil {
			if pool, err := database.DB.DB(); err == nil {
				if closeErr := pool.Close(); closeErr != nil {
					cleanupErrors = append(cleanupErrors, fmt.Errorf("close MySQL: %w", closeErr))
				}
			}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if closeErr := tracer.Shutdown(ctx); closeErr != nil {
			cleanupErrors = append(cleanupErrors, fmt.Errorf("flush traces: %w", closeErr))
		}
		err = errors.Join(append([]error{err}, cleanupErrors...)...)
	}()
	return rootCmd.Execute()
}
