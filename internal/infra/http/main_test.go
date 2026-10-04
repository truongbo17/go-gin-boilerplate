package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOutboundResponseLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", maxResponseBytes+1)))
	}))
	defer server.Close()
	client := &BaseRequest{Client: server.Client()}
	_, _, err := client.Get(context.Background(), server.URL, nil, nil)
	if !errors.Is(err, ErrResponseTooLarge) {
		t.Fatalf("got %v, want response size error", err)
	}
}
