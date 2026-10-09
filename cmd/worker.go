package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	appworker "github.com/truongbo17/go-gin-boilerplate/internal/app/worker"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/schedule"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/worker"
)

var StartWorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start background workers and schedules (requires Redis)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withRuntime(cmd, requirements{redis: true, tracing: true, worker: true}, func(r *runtime) error {
			registry, err := appworker.New(appworker.Dependencies{
				Config: r.config,
				DB:     r.db,
				Redis:  r.redis,
			})
			if err != nil {
				return err
			}
			scheduler, err := schedule.StartWith(r.redis, registry.Schedules)
			if err != nil {
				return err
			}
			server, err := worker.StartWith(r.config.Cache, r.log, registry.Handlers)
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
		})
	},
}
