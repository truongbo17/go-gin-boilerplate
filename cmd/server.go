package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/config"
	coreworker "github.com/truongbo17/go-gin-boilerplate/internal/app/core/worker"
	"github.com/truongbo17/go-gin-boilerplate/internal/health"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/worker/client"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares/limiter"
	"github.com/truongbo17/go-gin-boilerplate/internal/routes"
)

var StartServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withRuntime(cmd, requirements{database: true, cache: true, tracing: true}, func(r *runtime) error {
			return runServer(cmd, r)
		})
	},
}

func runServer(cmd *cobra.Command, r *runtime) error {
	if r.config.App.Env == config.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	rateLimiter, err := limiter.New(r.config.Cache.CacheStore, r.redis)
	if err != nil {
		return err
	}
	var jobs coreworker.Dispatcher
	if r.config.Mail.Enabled {
		queue := client.New(r.config.Cache)
		defer queue.Close()
		jobs = coreworker.New(queue)
	}
	router, err := routes.New(routes.Options{
		Config:         r.config,
		DB:             r.db,
		TokenBlacklist: r.cache,
		Readiness: health.Checker{
			DB:         r.db,
			Redis:      r.redis,
			CheckRedis: r.config.Cache.CacheStore == config.CacheStoreRedis,
		},
		Limiter: rateLimiter,
		Logger:  r.log,
		Jobs:    jobs,
	})
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:              ":" + r.config.App.Port,
		WriteTimeout:      30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Handler:           router,
	}

	serverError := make(chan error, 1)
	go func() { serverError <- server.ListenAndServe() }()
	r.log.Infof("Server is listening on port %s", r.config.App.Port)
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case err := <-serverError:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		r.log.Info("Server shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}
