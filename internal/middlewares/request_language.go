package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/utils"
)

func RequestLang() gin.HandlerFunc {
	return func(context *gin.Context) {
		language := context.GetHeader(config.HeaderLanguage)
		if !utils.CheckLanguage(language) {
			language = config.Vietnamese
		}
		context.Header(config.HeaderLanguage, language)
		context.Set(config.HeaderLanguage, language)

		//localize := i18n.NewLocalizer(i18n2.Bundle, language)
		//context.Set("localize", localize)

		context.Next()
	}
}
