package graphql

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	coregraphql "github.com/truongbo17/go-gin-boilerplate/internal/app/core/graphql"
	graphqlauth "github.com/truongbo17/go-gin-boilerplate/internal/app/graphql/auth"
	appLimiter "github.com/truongbo17/go-gin-boilerplate/internal/middlewares/limiter"
	authrepository "github.com/truongbo17/go-gin-boilerplate/internal/repository/auth"
	graphrepository "github.com/truongbo17/go-gin-boilerplate/internal/repository/graphql"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
	"github.com/ulule/limiter/v3"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB      *gorm.DB
	Auth    gin.HandlerFunc
	Limiter *appLimiter.Limiter
	Env     string
	Logger  *logrus.Logger
}

type Module struct {
	service coregraphql.Service
	handler http.Handler
	deps    Dependencies
}

func NewModule(deps Dependencies) *Module {
	service := coregraphql.Service{
		Store:       graphrepository.MetadataRepository{DB: deps.DB},
		Permissions: authrepository.NewPermissionRepository(deps.DB),
	}
	return &Module{service: service, handler: graphqlauth.NewHandler(deps.Env, service, deps.Logger), deps: deps}
}

func (module *Module) RegisterRoutes(api *gin.RouterGroup) {
	public := api.Group("/v1/public")
	public.GET("/enum-options/:key",
		module.deps.Limiter.Limit("graphql:public-enum", limiter.Rate{Period: time.Minute, Limit: 100}),
		module.publicEnum,
	)
	graph := api.Group("/v1/graphql")
	graph.POST("",
		module.deps.Limiter.Limit("graphql:query", limiter.Rate{Period: time.Minute, Limit: 100}),
		module.deps.Auth,
		module.execute,
	)
	graph.GET("/entity-options",
		module.deps.Limiter.Limit("graphql:entity-options", limiter.Rate{Period: time.Minute, Limit: 100}),
		module.deps.Auth,
		module.entityOptions,
	)
}

func (module *Module) publicEnum(c *gin.Context) {
	group, ok := coregraphql.EnumByKey(c.Param("key"))
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	response.ReturnSuccess(c, group, nil)
}

func (module *Module) execute(c *gin.Context) {
	if c.Request.ContentLength > 64<<10 {
		c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, 64<<10))
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	user := c.MustGet("user").(*models.User)
	request := c.Request.WithContext(graphqlauth.WithUser(c.Request.Context(), user))
	module.handler.ServeHTTP(c.Writer, request)
}

func (module *Module) entityOptions(c *gin.Context) {
	page, perPage := 1, 20
	if raw := c.Query("page"); raw != "" {
		var err error
		page, err = strconv.Atoi(raw)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	if raw := c.Query("per_page"); raw != "" {
		var err error
		perPage, err = strconv.Atoi(raw)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	user := c.MustGet("user").(*models.User)
	result, err := module.service.EntityOptions(c.Request.Context(), user.ID, c.Query("key"), c.Query("keyword"), page, perPage)
	if err != nil {
		switch {
		case errors.Is(err, coregraphql.ErrForbidden):
			c.AbortWithStatus(http.StatusForbidden)
		case errors.Is(err, coregraphql.ErrUnsupportedKey):
			c.AbortWithStatus(http.StatusBadRequest)
		case errors.Is(err, coregraphql.ErrInvalidInput):
			c.AbortWithStatus(http.StatusBadRequest)
		default:
			_ = c.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	response.ReturnSuccess(c, result, nil)
}
