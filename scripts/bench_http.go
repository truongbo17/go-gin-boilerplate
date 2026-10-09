package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

type workerResult struct {
	latencies []int64
	statuses  map[int]int
	failures  int
	bytes     int64
}

func main() {
	targetURL := flag.String("url", "http://127.0.0.1:8000/ping", "HTTP endpoint URL")
	concurrency := flag.Int("c", 16, "concurrent connections")
	duration := flag.Duration("duration", 10*time.Second, "measurement duration")
	warmup := flag.Duration("warmup", 2*time.Second, "warmup duration")
	flag.Parse()
	if *concurrency < 1 || *duration <= 0 || *warmup < 0 {
		fmt.Fprintln(os.Stderr, "concurrency and duration must be positive; warmup cannot be negative")
		os.Exit(2)
	}
	transport := &http.Transport{
		MaxIdleConns:          *concurrency,
		MaxIdleConnsPerHost:   *concurrency,
		MaxConnsPerHost:       *concurrency,
		DialContext:           (&net.Dialer{Timeout: 2 * time.Second}).DialContext,
		ResponseHeaderTimeout: 5 * time.Second,
		IdleConnTimeout:       30 * time.Second,
		DisableCompression:    true,
	}
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport, Timeout: 5 * time.Second}
	url := *targetURL
	if _, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
		fmt.Fprintln(os.Stderr, "invalid URL:", err)
		os.Exit(2)
	}
	if *warmup > 0 {
		run(client, url, *concurrency, *warmup, false)
	}
	results, elapsed := run(client, url, *concurrency, *duration, true)
	var latencies []int64
	statuses := make(map[int]int)
	failures := 0
	var bytes int64
	for _, result := range results {
		latencies = append(latencies, result.latencies...)
		for code, count := range result.statuses {
			statuses[code] += count
		}
		failures += result.failures
		bytes += result.bytes
	}
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	percentile := func(p float64) float64 {
		if len(latencies) == 0 {
			return 0
		}
		i := int(float64(len(latencies)-1) * p)
		return float64(latencies[i]) / 1000
	}
	total := len(latencies)
	fmt.Printf("url=%s concurrency=%d elapsed=%.3fs total=%d rps=%.1f p50=%.1fus p95=%.1fus p99=%.1fus max=%.1fus statuses=%v errors=%d body_bytes=%d\n", url, *concurrency, elapsed.Seconds(), total, float64(total)/elapsed.Seconds(), percentile(.5), percentile(.95), percentile(.99), percentile(1), statuses, failures, bytes)
}

func run(client *http.Client, url string, concurrency int, duration time.Duration, record bool) ([]workerResult, time.Duration) {
	results := make([]workerResult, concurrency)
	var wg sync.WaitGroup
	start := make(chan struct{})
	deadline := time.Now().Add(duration)
	for i := range results {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			result := &results[index]
			result.statuses = make(map[int]int)
			<-start
			for time.Now().Before(deadline) {
				begun := time.Now()
				request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
				response, err := client.Do(request)
				if err != nil {
					result.failures++
				} else {
					n, readErr := io.Copy(io.Discard, response.Body)
					closeErr := response.Body.Close()
					if readErr != nil || closeErr != nil {
						result.failures++
					}
					result.bytes += n
					result.statuses[response.StatusCode]++
				}
				if record {
					result.latencies = append(result.latencies, time.Since(begun).Nanoseconds())
				}
			}
		}(i)
	}
	started := time.Now()
	close(start)
	wg.Wait()
	return results, time.Since(started)
}
