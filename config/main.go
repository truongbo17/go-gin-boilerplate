package config

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	App      `mapstructure:",squash"`
	Cors     `mapstructure:",squash"`
	Database `mapstructure:",squash"`
	Cache    `mapstructure:",squash"`
	Auth     `mapstructure:",squash"`
	Tracer   `mapstructure:",squash"`
}

const (
	// DebugMode app env debug stg.
	DebugMode string = "debug"
	// ReleaseMode app env debug production.
	ReleaseMode string = "release"
	// LocalMode for dev.
	LocalMode string = "local"
)

func (config *Config) validate() error {
	err := validation.ValidateStruct(config,
		// App
		validation.Field(&config.App.Port, is.Port),
		validation.Field(&config.App.Env, validation.In(DebugMode, ReleaseMode, LocalMode)),

		// CORS
		validation.Field(&config.Cors.AllowOrigin),

		// Database
		validation.Field(&config.Database.Master.Port, validation.Required, is.Port),
		validation.Field(&config.Database.Master.Host, validation.Required, is.Host),
		validation.Field(&config.Database.Master.Username, validation.Required),
		validation.Field(&config.Database.Master.Database, validation.Required),

		// Cache
		validation.Field(&config.Cache.CacheStore, validation.In(CacheStoreLocal, CacheStoreRedis)),

		// Auth
		validation.Field(&config.Auth.JWTSecretKey, validation.Required, validation.Length(32, 0)),
		validation.Field(&config.Tracer.SampleRatio, validation.Min(0.0), validation.Max(1.0)),
	)
	if err != nil {
		return err
	}
	if config.Cache.CacheStore == CacheStoreRedis {
		return validation.ValidateStruct(config,
			validation.Field(&config.Cache.RedisHost, validation.Required, is.Host),
			validation.Field(&config.Cache.RedisPort, validation.Required, is.Port),
		)
	}
	return nil
}

var EnvConfig *Config

func setupConfig() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.SetDefault("APP_ENV", "debug")
	viper.SetDefault("APP_PORT", "8000")
	viper.SetDefault("APP_TRUSTED_PROXIES", "")
	viper.SetDefault("CORS_ALLOW_ORIGIN", "http://localhost:3000")
	viper.SetDefault("CACHE_STORE", "local")

	viper.SetDefault("JWT_ACCESS_EXPIRATION_MINUTES", 24*60)
	viper.SetDefault("JWT_REFRESH_EXPIRATION_DAYS", 30)

	viper.SetDefault("TRACER_ENABLE", "false")
	viper.SetDefault("TRACER_SAMPLE_RATIO", 0.1)
	viper.SetDefault("ELASTIC_APM_ENVIRONMENT", "staging")
	for _, key := range []string{
		"APP_NAME", "APP_URL", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_DATABASE",
		"REDIS_HOST", "REDIS_PORT", "REDIS_USERNAME", "REDIS_PASSWORD", "JWT_SECRET",
		"ELASTIC_APM_SERVER_URL", "ELASTIC_APM_SECRET_TOKEN", "ELASTIC_APM_SERVICE_NAME",
		"ELASTIC_APM_GLOBAL_LABELS",
	} {
		if err := viper.BindEnv(key); err != nil {
			panic(err)
		}
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	if err := config.validate(); err != nil {
		panic(err)
	}

	return config
}

func Init() {
	EnvConfig = setupConfig()

	fmt.Println("Success init config")
}
