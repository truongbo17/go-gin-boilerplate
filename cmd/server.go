package cmd

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/limiter"
	"github.com/truongbo17/go-gin-boilerplate/internal/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	StartServerCmd = &cobra.Command{
		Use:   "server",
		Short: `Start the server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return start()
		},
	}
)

func start() error {
	EnvConfig := config.EnvConfig
	storeCache := EnvConfig.Cache.CacheStore

	limiter.InitLimiterStore(storeCache)

	routes.Init()
	r := routes.Router

	server := &http.Server{
		Addr:              ":" + EnvConfig.App.Port,
		WriteTimeout:      time.Second * 30,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       time.Second * 30,
		MaxHeaderBytes:    1 << 20,
		Handler:           r,
	}

	log.Printf("Server is now listening at port: %s. Good luck!", EnvConfig.App.Port)

	serverError := make(chan error, 1)
	go func() { serverError <- server.ListenAndServe() }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)
	select {
	case err := <-serverError:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case i := <-quit:
		log.Println("Server receive a signal: ", i.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("Server exiting.")
	return nil
}
