package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestRedisRequiredOnlyWhenSelected(t *testing.T) {
	config := &Config{
		App:      App{Port: "8000", Env: LocalMode},
		Database: Database{Master: Master{Host: "127.0.0.1", Port: "3306", Username: "test", Database: "test", MaxOpenConns: 30, MaxIdleConns: 15}},
		Cache:    Cache{CacheStore: CacheStoreLocal},
		Auth:     Auth{JWTSecretKey: strings.Repeat("x", 32)},
	}
	if err := config.validate(); err != nil {
		t.Fatal(err)
	}
	config.Cache.CacheStore = CacheStoreRedis
	if err := config.validate(); err == nil {
		t.Fatal("Redis configuration was not required")
	}
	config.Cache.RedisHost = "127.0.0.1"
	config.Cache.RedisPort = "6379"
	if err := config.validate(); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentOnlyConfiguration(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv("APP_ENV", "local")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_USER", "test")
	t.Setenv("DB_DATABASE", "test")
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
	config := setupConfig()
	if config.Database.Master.Username != "test" || config.App.Port != "8000" || config.Database.Master.MaxOpenConns != 30 {
		t.Fatalf("environment-only configuration was not loaded: %+v", config)
	}
}

func TestDatabasePoolConfigurationIsValidated(t *testing.T) {
	config := &Config{
		App:      App{Port: "8000", Env: LocalMode},
		Database: Database{Master: Master{Host: "127.0.0.1", Port: "3306", Username: "test", Database: "test", MaxOpenConns: 10, MaxIdleConns: 11}},
		Cache:    Cache{CacheStore: CacheStoreLocal},
		Auth:     Auth{JWTSecretKey: strings.Repeat("x", 32)},
	}
	if err := config.validate(); err == nil {
		t.Fatal("accepted more idle than open database connections")
	}
}

func TestTraceSampleRatioIsBounded(t *testing.T) {
	config := &Config{
		App:      App{Port: "8000", Env: LocalMode},
		Database: Database{Master: Master{Host: "127.0.0.1", Port: "3306", Username: "test", Database: "test", MaxOpenConns: 30, MaxIdleConns: 15}},
		Cache:    Cache{CacheStore: CacheStoreLocal},
		Auth:     Auth{JWTSecretKey: strings.Repeat("x", 32)},
		Tracer:   Tracer{SampleRatio: 1.1},
	}
	if err := config.validate(); err == nil {
		t.Fatal("accepted sampling ratio above one")
	}
}
