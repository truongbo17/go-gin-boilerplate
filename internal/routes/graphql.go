package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	graphqlauth "github.com/truongbo17/go-gin-boilerplate/internal/app/graphql/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/v1/auth"
	"github.com/ulule/limiter/v3"
)

func registerGraphQL(r *gin.RouterGroup, module *auth.Module, options Options) {
	handler := graphqlauth.NewHandler(options.Config.App.Env)
	r.POST("/v1/graphql",
		options.Limiter.Limit("graphql:query", limiter.Rate{Period: time.Minute, Limit: 100}),
		module.AuthMiddleware(),
		func(c *gin.Context) {
			user := c.MustGet("user").(*models.User)
			request := c.Request.WithContext(graphqlauth.WithUser(c.Request.Context(), user))
			handler.ServeHTTP(c.Writer, request)
		},
	)
}
