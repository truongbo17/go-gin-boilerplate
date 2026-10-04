package cmd

import (
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
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "ggb",
	Short: "Reusable Gin API base",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "version" {
			return nil
		}
		config.Init()
		if config.EnvConfig.App.Env == config.ReleaseMode {
			gin.SetMode(gin.ReleaseMode)
		}
		logger.Init()
		redis.ConnectRedis()
		cache.InitCache()
		database.ConnectDatabase()
		httpclient.InitBaseRequest()
		tracer.InitTracerOTEL()
		if config.EnvConfig.Cache.CacheStore == config.CacheStoreRedis {
			client.InitClient()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(StartServerCmd, StartWorkerCmd, cli.VersionCmd, cli.MigrateCmd, cli.CreateUserCmd)
}

func Execute() {
	defer func() {
		if client.WorkerClient != nil {
			_ = client.WorkerClient.Close()
		}
		tracer.DownAMPTracerProvider()
	}()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
