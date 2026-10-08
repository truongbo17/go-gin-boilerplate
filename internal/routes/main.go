package routes

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var Router *gin.Engine

func Init() {
	Router = gin.New()
	var trustedProxies []string
	for _, proxy := range strings.Split(config.EnvConfig.App.TrustedProxies, ",") {
		if proxy = strings.TrimSpace(proxy); proxy != "" {
			trustedProxies = append(trustedProxies, proxy)
		}
	}
	if err := Router.SetTrustedProxies(trustedProxies); err != nil {
		panic(fmt.Errorf("invalid APP_TRUSTED_PROXIES: %w", err))
	}
	Router.Use(gin.Recovery())
	if config.EnvConfig.Tracer.Enable {
		Router.Use(otelgin.Middleware(config.EnvConfig.Tracer.APM.ServiceName))
	}
	Router.Use(middlewares.SecurityHeaders(), middlewares.LimitRequestBody(), middlewares.RequestID(), middlewares.RequestLang(), middlewares.RequestLogger(), middlewares.Cors())
	LoadPublic(Router)
	api := Router.Group("/api")
	auth.LoadAuthV1(api)
	auth.LoadRBACV1(api)
}
