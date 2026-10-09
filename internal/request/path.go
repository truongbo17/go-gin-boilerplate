package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/internal/utils"
)

// PathID accepts only positive IDs representable by the 32-bit unsigned integers.
func PathID(ctx *gin.Context) (uint, bool) {
	id, err := utils.ParseID(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return 0, false
	}
	return id, true
}
