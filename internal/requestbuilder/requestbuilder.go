package requestbuilder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/rehttp"
)

const (
	DefaultMaxRetries         = 5
	DefaultTimeout            = 10 * time.Second
	DefaultBaseExpJitterDelay = 100 * time.Millisecond
	DefaultMaxExpJitterDelay  = 5 * time.Second
)

type RequestBuilder struct {
	method      string
	url         string
	input       any
	headers     http.Header
	httpClient  *http.Client
	queryParams map[string][]string
}

func New(method, url string) *RequestBuilder {
	return &RequestBuilder{
		method: method,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		headers:     make(http.Header),
		url:         url,
		queryParams: make(map[string][]string),
		input:       nil,
	}
}

func (r *RequestBuilder) WithClient(httpClient *http.Client) *RequestBuilder {
	r.httpClient = httpClient
	return r
}

func (r *RequestBuilder) WithTimeout(timeout time.Duration) *RequestBuilder {
	r.httpClient.Timeout = timeout
	return r
}

func (r *RequestBuilder) WithBody(input any) *RequestBuilder {
	r.headers.Set("Content-Type", "application/json")
	r.input = input
	return r
}

func (r *RequestBuilder) WithPathParam(key string, value string) *RequestBuilder {
	r.url = strings.ReplaceAll(r.url, fmt.Sprintf(":%s", key), value)
	return r
}

func (r *RequestBuilder) WithQueryParam(key string, values ...string) *RequestBuilder {
	r.queryParams[key] = values
	return r
}

func (r *RequestBuilder) WithHeader(key string, values ...string) *RequestBuilder {
	for _, v := range values {
		r.headers.Add(key, v)
	}

	return r
}

func (r *RequestBuilder) WithRetry(additionalStatuses ...int) *RequestBuilder {
	retryErrorsAndStatuses := []rehttp.RetryFn{
		rehttp.RetryTimeoutErr(),
		rehttp.RetryStatusInterval(500, 600),
	}

	if len(additionalStatuses) > 0 {
		retryErrorsAndStatuses = append(retryErrorsAndStatuses, rehttp.RetryStatuses(additionalStatuses...))
	}

	r.httpClient.Transport = rehttp.NewTransport(
		r.httpClient.Transport,
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(DefaultMaxRetries),
			rehttp.RetryAny(retryErrorsAndStatuses...),
		),
		rehttp.ExpJitterDelay(DefaultBaseExpJitterDelay, DefaultMaxExpJitterDelay),
	)

	return r
}

func (r *RequestBuilder) WithCustomRetry(retryFn rehttp.RetryFn, delayFn rehttp.DelayFn) *RequestBuilder {
	r.httpClient.Transport = rehttp.NewTransport(r.httpClient.Transport, retryFn, delayFn)

	return r
}

func (r *RequestBuilder) Build(ctx context.Context) (*http.Request, error) {
	var body io.Reader = nil
	if r.input != nil {
		marshalled, err := json.Marshal(r.input)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input: %w", err)
		}

		body = bytes.NewReader(marshalled)
	}

	req, err := http.NewRequestWithContext(ctx, r.method, r.url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range r.queryParams {
		q := req.URL.Query()
		for _, value := range v {
			q.Add(k, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header = r.headers
	if ctx.Value("correlation_id") != nil {
		req.Header.Set("X-Correlation-ID", ctx.Value("correlation_id").(string))
	}

	return req, nil
}

func (r *RequestBuilder) Do(ctx context.Context) (*http.Response, error) {
	req, err := r.Build(ctx)
	if err != nil {
		return nil, err
	}

	return r.httpClient.Do(req)
}
