package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"gorm.io/gorm"
)

var benchmarkService AuthService
var benchmarkToken string
var benchmarkDB *gorm.DB

//go:noinline
func constructAuthServiceForBenchmark() AuthService {
	return NewAuthService()
}

func benchmarkUser(b *testing.B) *models.User {
	b.Helper()
	previous := config.EnvConfig
	config.EnvConfig = &config.Config{Auth: config.Auth{
		JWTSecretKey:               strings.Repeat("x", 32),
		JWTAccessExpirationMinutes: 60,
	}}
	b.Cleanup(func() { config.EnvConfig = previous })
	user := &models.User{Username: "benchmark", Email: "benchmark@example.com", Password: "hashed-password"}
	user.ID = 1
	return user
}

func BenchmarkAuthServiceConstruction(b *testing.B) {
	previousDB := database.DB
	database.DB = &gorm.DB{}
	b.Cleanup(func() { database.DB = previousDB })
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkService = constructAuthServiceForBenchmark()
		benchmarkDB = benchmarkService.UserRepository.DB
	}
}

func BenchmarkGenerateTokenWithServiceConstruction(b *testing.B) {
	user := benchmarkUser(b)
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service := NewAuthService()
		token, err := service.GenerateToken(ctx, enums.TokenTypeAccess, user)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkToken = token
	}
}

func BenchmarkTokenRoute(b *testing.B) {
	user := benchmarkUser(b)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/token", func(c *gin.Context) {
		service := NewAuthService()
		token, err := service.GenerateToken(c.Request.Context(), enums.TokenTypeAccess, user)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.String(http.StatusOK, "%s", token)
	})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/token", nil))
		if response.Code != http.StatusOK {
			b.Fatalf("route returned %d", response.Code)
		}
		benchmarkToken = response.Body.String()
	}
}
