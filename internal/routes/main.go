package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var Router *gin.Engine

func Init() {
	Router = gin.New()
	_ = Router.SetTrustedProxies(nil)
	Router.Use(gin.Recovery())
	if config.EnvConfig.Tracer.Enable {
		Router.Use(otelgin.Middleware(config.EnvConfig.Tracer.APM.ServiceName))
	}
	Router.Use(middlewares.RequestID(), middlewares.RequestLang(), middlewares.RequestLogger(), middlewares.Cors(), middlewares.RateInternalLimit())
	LoadPublic(Router)
	api := Router.Group("/api")
	auth.LoadAuthV1(api)
	auth.LoadRBACV1(api)
}
