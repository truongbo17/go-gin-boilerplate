package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
)

func ReturnPasswordResetRequested(ctx *gin.Context) {
	response.ReturnSuccessStatus(ctx, http.StatusAccepted,
		authMessage(ctx.GetString(config.HeaderLanguage), "auth.reset_requested"), nil, nil)
}
