package config

type Tracer struct {
	Enable bool `mapstructure:"TRACER_ENABLE"`
	APM    `mapstructure:",squash"`
}

type APM struct {
	Url         string `mapstructure:"ELASTIC_APM_SERVER_URL"`
	SecretKey   string `mapstructure:"ELASTIC_APM_SECRET_TOKEN"`
	ServiceName string `mapstructure:"ELASTIC_APM_SERVICE_NAME"`
	Environment string `mapstructure:"ELASTIC_APM_ENVIRONMENT"`
	Labels      string `mapstructure:"ELASTIC_APM_GLOBAL_LABELS"`
}
