package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/i18n"
)

func RequestLang() gin.HandlerFunc {
	return func(context *gin.Context) {
		language := context.GetHeader(config.HeaderLanguage)
		if !i18n.Supported(language) {
			language = config.Vietnamese
		}
		context.Header(config.HeaderLanguage, language)
		context.Set(config.HeaderLanguage, language)

		context.Next()
	}
}
