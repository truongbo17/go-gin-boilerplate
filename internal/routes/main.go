package routes

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/services"
	coreworker "github.com/truongbo17/go-gin-boilerplate/internal/app/core/worker"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/health"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares/limiter"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/gorm"
)

type Options struct {
	Config         config.Config
	DB             *gorm.DB
	TokenBlacklist services.TokenBlacklist
	Readiness      health.Checker
	Limiter        *limiter.Limiter
	Logger         *logrus.Logger
	Jobs           coreworker.Dispatcher
}

func New(options Options) (*gin.Engine, error) {
	if options.DB == nil || options.TokenBlacklist == nil || options.Limiter == nil || options.Logger == nil {
		return nil, fmt.Errorf("HTTP router dependencies are incomplete")
	}
	authModule := auth.NewModule(auth.Dependencies{
		DB:             options.DB,
		TokenBlacklist: options.TokenBlacklist,
		Auth:           options.Config.Auth,
		Limiter:        options.Limiter,
		Jobs:           options.Jobs,
	})

	router := gin.New()
	var trustedProxies []string
	for proxy := range strings.SplitSeq(options.Config.App.TrustedProxies, ",") {
		if proxy = strings.TrimSpace(proxy); proxy != "" {
			trustedProxies = append(trustedProxies, proxy)
		}
	}
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		return nil, fmt.Errorf("invalid APP_TRUSTED_PROXIES: %w", err)
	}
	router.Use(gin.Recovery())
	if options.Config.Tracer.Enable {
		router.Use(otelgin.Middleware(options.Config.Tracer.APM.ServiceName))
	}
	router.Use(middlewares.SecurityHeaders(), middlewares.LimitRequestBody(), middlewares.RequestID(), middlewares.RequestLang(), middlewares.RequestLoggerWith(options.Logger), middlewares.CorsWith(options.Config.Cors.AllowOrigin))
	router.Use(func(c *gin.Context) { c.Set("app_url", options.Config.App.Url); c.Next() })
	LoadPublic(router, options.Readiness)
	registerOpenAPI(router)
	api := router.Group("/api")
	authModule.RegisterRoutes(api)
	registerGraphQL(api, authModule, options)
	return router, nil
}
