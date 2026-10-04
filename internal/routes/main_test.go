package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/limiter"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/logger"
)

func TestBaseRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.EnvConfig = &config.Config{App: config.App{Env: config.LocalMode, Name: "test"}, Cors: config.Cors{AllowOrigin: "http://localhost:3000"}, Cache: config.Cache{CacheStore: config.CacheStoreLocal}}
	logger.LogrusLogger = logrus.New()
	limiter.InitLimiterStore(config.CacheStoreLocal)
	Init()

	cases := []struct {
		path string
		want int
	}{
		{"/ping", http.StatusOK},
		{"/api/v1/auth/me", http.StatusUnauthorized},
		{"/api/v1/campaign", http.StatusNotFound},
	}
	for _, tc := range cases {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, tc.path, nil)
		Router.ServeHTTP(recorder, request)
		if recorder.Code != tc.want {
			t.Errorf("%s: got %d, want %d", tc.path, recorder.Code, tc.want)
		}
	}
}
