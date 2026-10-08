package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/schedule"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/worker"
)

var StartWorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start background workers and schedules (requires Redis)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.EnvConfig.Cache.CacheStore != config.CacheStoreRedis {
			return fmt.Errorf("worker requires CACHE_STORE=redis")
		}
		scheduler, err := schedule.Start()
		if err != nil {
			return err
		}
		server, err := worker.Start()
		if err != nil {
			_ = scheduler.Shutdown()
			return err
		}
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		<-ctx.Done()
		if err := scheduler.Shutdown(); err != nil {
			server.Shutdown()
			return fmt.Errorf("stop scheduler: %w", err)
		}
		server.Shutdown()
		return nil
	},
}
