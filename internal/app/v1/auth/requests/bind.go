package requests

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func bindRequest(c *gin.Context, target any) bool {
	if err := c.ShouldBind(target); err != nil {
		var limitError *http.MaxBytesError
		if errors.As(err, &limitError) {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		} else {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
		}
		return false
	}
	return true
}
