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
		{"/ready", http.StatusServiceUnavailable},
		{"/api/v1/auth/me", http.StatusUnauthorized},
		{"/api/v1/unknown", http.StatusNotFound},
	}
	for _, tc := range cases {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, tc.path, nil)
		Router.ServeHTTP(recorder, request)
		if recorder.Code != tc.want {
			t.Errorf("%s: got %d, want %d", tc.path, recorder.Code, tc.want)
		}
		if recorder.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Errorf("%s: missing security header", tc.path)
		}
	}
}

func TestLoginHasDedicatedRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.EnvConfig = &config.Config{App: config.App{Env: config.LocalMode}, Cors: config.Cors{AllowOrigin: "http://localhost:3000"}, Cache: config.Cache{CacheStore: config.CacheStoreLocal}}
	logger.LogrusLogger = logrus.New()
	limiter.InitLimiterStore(config.CacheStoreLocal)
	Init()
	for i := 0; i < 11; i++ {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		Router.ServeHTTP(recorder, request)
		if i == 10 && recorder.Code != http.StatusTooManyRequests {
			t.Fatalf("request %d: got %d, want 429", i+1, recorder.Code)
		}
	}
}

func TestLargeRequestIsRejected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.EnvConfig = &config.Config{App: config.App{Env: config.LocalMode}, Cors: config.Cors{AllowOrigin: "http://localhost:3000"}, Cache: config.Cache{CacheStore: config.CacheStoreLocal}}
	logger.LogrusLogger = logrus.New()
	limiter.InitLimiterStore(config.CacheStoreLocal)
	Init()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	request.ContentLength = 1<<20 + 1
	Router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("got %d, want 413", recorder.Code)
	}
}
