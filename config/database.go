package config

type Database struct {
	Master `mapstructure:",squash"`
}

const (
	DefaultStringSizeMySql int    = 256
	TableMigrate           string = "migrations"
)

type Master struct {
	Username               string `mapstructure:"DB_USER"`
	Password               string `mapstructure:"DB_PASS"`
	Host                   string `mapstructure:"DB_HOST"`
	Port                   string `mapstructure:"DB_PORT"`
	Database               string `mapstructure:"DB_DATABASE"`
	MaxOpenConns           int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns           int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetimeMinutes int    `mapstructure:"DB_CONN_MAX_LIFETIME_MINUTES"`
	ConnMaxIdleTimeMinutes int    `mapstructure:"DB_CONN_MAX_IDLE_TIME_MINUTES"`
}
