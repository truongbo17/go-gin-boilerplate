package http

import (
	"context"
	"errors"
	"io"
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

func TestOutboundResponseBodyRemainsReadable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, "accepted")
	}))
	defer server.Close()
	client := &BaseRequest{Client: server.Client()}
	response, body, err := client.Get(context.Background(), server.URL, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	readBack, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusAccepted || string(body) != "accepted" || string(readBack) != "accepted" {
		t.Fatalf("response status=%d, body=%q, readable=%q", response.StatusCode, body, readBack)
	}
}
