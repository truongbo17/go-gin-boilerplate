package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestListUserPageSizeIsBounded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/users?per_page=1000000", nil)
	ListUserValidator()(ctx)
	value, exists := ctx.Get("ListUserRequest")
	if !exists || value.(ListUserRequest).PerPage != 100 {
		t.Fatalf("unexpected page size: %v", value)
	}
}

func TestListUserRejectsExtremePage(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/users?page=10001", nil)
	ListUserValidator()(ctx)
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("got %d, want 422", recorder.Code)
	}
}
