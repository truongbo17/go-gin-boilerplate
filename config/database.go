package config

type Database struct {
	Master `mapstructure:",squash"`
}

const (
	DefaultStringSizeMySql int    = 256
	TableMigrate           string = "migrations"
)

type Master struct {
	Username string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASS"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Database string `mapstructure:"DB_DATABASE"`
}
