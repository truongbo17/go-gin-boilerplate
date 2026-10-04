package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseID accepts only positive IDs representable by the MySQL UNSIGNED INT columns.
func parseID(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return 0, false
	}
	return uint(id), true
}
