package request

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Bind(c *gin.Context, target any) bool {
	if err := c.ShouldBind(target); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		} else {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
		}
		return false
	}
	return true
}
