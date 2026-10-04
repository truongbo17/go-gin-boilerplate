package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBusinessErrorsUseHTTPErrorStatus(t *testing.T) {
	for _, test := range []struct{ code, want int }{
		{ErrAuthWrongPassword, http.StatusUnauthorized},
		{ErrForbidden, http.StatusForbidden},
		{ErrRoleNotFound, http.StatusNotFound},
		{ErrAuthUserExists, http.StatusConflict},
		{ErrRoleInternalError, http.StatusInternalServerError},
	} {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		ReturnError(ctx, http.StatusOK, test.code, nil)
		if recorder.Code != test.want {
			t.Errorf("code %d: got HTTP %d, want %d", test.code, recorder.Code, test.want)
		}
	}
}
