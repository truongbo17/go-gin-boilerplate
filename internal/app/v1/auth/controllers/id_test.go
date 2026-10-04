package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParseIDRejectsInvalidAndOverflow(t *testing.T) {
	for _, value := range []string{"", "0", "not-a-number", "4294967296"} {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = gin.Params{{Key: "id", Value: value}}
		_, ok := parseID(ctx)
		if ok || ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("accepted ID %q", value)
		}
	}
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Params = gin.Params{{Key: "id", Value: "4294967295"}}
	id, ok := parseID(ctx)
	if !ok || id != 4294967295 {
		t.Fatalf("valid ID rejected: %v, %v", id, ok)
	}
}
