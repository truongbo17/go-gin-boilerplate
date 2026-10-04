package http

import (
	"bytes"
	"context"
	json "github.com/json-iterator/go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"net/http"
	"net/url"
	"time"
)

var Request *BaseRequest

func InitBaseRequest() {
	Request = &BaseRequest{
		Client: &http.Client{
			Timeout:   20 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

type BaseRequest struct {
	Client *http.Client
}

func (r *BaseRequest) Do(ctx context.Context, method, urlStr string, headers map[string]string, body interface{}) (*http.Response, []byte, error) {
	var requestBody []byte
	var err error
	if body != nil {
		requestBody, err = json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, respBody, nil
}

func (r *BaseRequest) Get(ctx context.Context, baseURL string, headers map[string]string, queryParams map[string]string) (*http.Response, []byte, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, err
	}

	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return r.Do(ctx, http.MethodGet, u.String(), headers, nil)
}

func (r *BaseRequest) Post(ctx context.Context, url string, headers map[string]string, body interface{}) (*http.Response, []byte, error) {
	return r.Do(ctx, http.MethodPost, url, headers, body)
}

func (r *BaseRequest) Put(ctx context.Context, url string, headers map[string]string, body interface{}) (*http.Response, []byte, error) {
	return r.Do(ctx, http.MethodPut, url, headers, body)
}

func (r *BaseRequest) Delete(ctx context.Context, url string, headers map[string]string, body interface{}) (*http.Response, []byte, error) {
	return r.Do(ctx, http.MethodDelete, url, headers, body)
}
