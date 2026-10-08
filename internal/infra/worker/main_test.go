package worker

import (
	"errors"
	"testing"
	"time"
)

func TestRetryDelayIsBounded(t *testing.T) {
	for _, tc := range []struct {
		attempt int
		want    time.Duration
	}{
		{0, 5 * time.Second},
		{1, 10 * time.Second},
		{2, 20 * time.Second},
		{100, 5 * time.Minute},
	} {
		if got := retryDelay(tc.attempt, errors.New("temporary"), nil); got != tc.want {
			t.Fatalf("attempt %d: got %v, want %v", tc.attempt, got, tc.want)
		}
	}
	custom := &RateLimitError{RetryInDuration: 15 * time.Second}
	if got := retryDelay(100, custom, nil); got != custom.RetryInDuration {
		t.Fatalf("custom retry delay = %v", got)
	}
}
