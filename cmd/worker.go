package cmd

import (
	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/schedule"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/worker"
)

var StartWorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start background workers and schedules (requires Redis)",
	Run: func(cmd *cobra.Command, args []string) {
		if config.EnvConfig.Cache.CacheStore != config.CacheStoreRedis {
			panic("worker requires CACHE_STORE=redis")
		}
		schedule.Init()
		worker.InitServer()
	},
}
