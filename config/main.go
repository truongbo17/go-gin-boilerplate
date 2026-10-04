package config

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/viper"
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
	return validation.ValidateStruct(config,
		// App
		validation.Field(&config.App.Port, is.Port),
		validation.Field(&config.App.Env, validation.In(DebugMode, ReleaseMode, LocalMode)),

		// CORS
		validation.Field(&config.Cors.AllowOrigin),

		// Database
		validation.Field(&config.Database.Master.Port, is.Port),
		validation.Field(&config.Database.Master.Host, is.Host),

		// Cache
		validation.Field(&config.Cache.CacheStore, validation.In(CacheStoreLocal, CacheStoreRedis)),

		// Redis
		validation.Field(&config.Cache.RedisPort, is.Port),
		validation.Field(&config.Cache.RedisHost, is.Host),

		// Auth
		validation.Field(&config.Auth.JWTSecretKey, validation.Required, validation.Length(32, 0)),
	)
}

var EnvConfig *Config

func setupConfig() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.SetDefault("APP_ENV", "debug")
	viper.SetDefault("APP_PORT", "8000")
	viper.SetDefault("CORS_ALLOW_ORIGIN", "http://localhost:3000")
	viper.SetDefault("CACHE_STORE", "local")

	viper.SetDefault("JWT_ACCESS_EXPIRATION_MINUTES", 24*60)
	viper.SetDefault("JWT_REFRESH_EXPIRATION_DAYS", 30)

	viper.SetDefault("TRACER_ENABLE", "false")
	viper.SetDefault("ELASTIC_APM_ENVIRONMENT", "staging")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
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
