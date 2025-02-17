package config

const (
	PathLog                       string = "storage/logs/%s.log"
	DefaultScheduleLockRedisRetry int    = 1
	HeaderRequestID               string = "X-Request-ID"
)

type AppConfig struct {
	Env      string `mapstructure:"APP_ENV"`
	Name     string `mapstructure:"APP_NAME"`
	Port     string `mapstructure:"APP_PORT"`
	IsWorker bool   `mapstructure:"APP_WORKER"`
}
