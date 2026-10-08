package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	ulimiter "github.com/ulule/limiter/v3"
)

func TestPoliciesHaveIndependentCounters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	InitLimiterStore(config.CacheStoreLocal)
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.GET("/login", Limit("login", ulimiter.Rate{Period: time.Minute, Limit: 1}), func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/api", Limit("api", ulimiter.Rate{Period: time.Minute, Limit: 1}), func(c *gin.Context) { c.Status(http.StatusOK) })

	for _, tc := range []struct {
		path string
		want int
	}{
		{"/login", http.StatusOK},
		{"/api", http.StatusOK},
		{"/login", http.StatusTooManyRequests},
		{"/api", http.StatusTooManyRequests},
	} {
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if recorder.Code != tc.want {
			t.Fatalf("%s: got %d, want %d", tc.path, recorder.Code, tc.want)
		}
		if got := recorder.Header().Get("X-RateLimit-Limit"); got != "1" {
			t.Fatalf("%s: limit header = %q", tc.path, got)
		}
		if tc.want == http.StatusTooManyRequests && recorder.Header().Get("Retry-After") == "" {
			t.Fatalf("%s: missing Retry-After header", tc.path)
		}
	}
}
