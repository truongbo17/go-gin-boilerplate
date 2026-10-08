package config

const (
	PathLog                       = "storage/logs/%s.log"
	DefaultScheduleLockRedisRetry = 1
	HeaderRequestID               = "X-Request-ID"
	HeaderLanguage                = "X-Language"
	PathV1                        = "v1"
	English                       = "en"
	Vietnamese                    = "vi"
)

var LanguageAvailable = [2]string{English, Vietnamese}

type App struct {
	Env            string `mapstructure:"APP_ENV"`
	Url            string `mapstructure:"APP_URL"`
	Name           string `mapstructure:"APP_NAME"`
	Port           string `mapstructure:"APP_PORT"`
	TrustedProxies string `mapstructure:"APP_TRUSTED_PROXIES"`
}
