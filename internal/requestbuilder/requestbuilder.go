package requestbuilder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/rehttp"
	"github.com/cybre/home-inventory/internal/cache"
	"github.com/cybre/home-inventory/internal/utils"
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
	r.headers[key] = values

	return r
}

func (r *RequestBuilder) WithCache(cache *cache.Cache[string], key string) *RequestBuilder {
	if !methodIsCacheable(r.method) {
		panic("WithCache can only be used with GET and HEAD requests")
	}

	r.httpClient.Transport = NewAutoCacheTransport(r.httpClient.Transport, cache, key)

	return r
}

func (r *RequestBuilder) WithSetCacheFn(cache *cache.Cache[string], setCacheFns ...SetCacheFn) *RequestBuilder {
	if !methodIsAction(r.method) {
		panic("WithSetCacheFn can only be used with POST, PUT, PATCH, and DELETE requests")
	}

	r.httpClient.Transport = NewSetCacheTransport(r.httpClient.Transport, cache, setCacheFns...)

	return r
}

func (r *RequestBuilder) WithInvalidateCache(cache *cache.Cache[string], keys ...string) *RequestBuilder {
	setCacheFns := utils.Map(keys, func(_ uint, key string) SetCacheFn {
		return func() (string, any) {
			return key, nil
		}
	})

	return r.WithSetCacheFn(cache, setCacheFns...)
}

func (r *RequestBuilder) WithRetry(additionalStatuses ...int) *RequestBuilder {
	retryErrorsAndStatuses := []rehttp.RetryFn{
		rehttp.RetryTimeoutErr(),
		rehttp.RetryStatusInterval(500, 600),
		rehttp.RetryIsErr(func(err error) bool {
			// Retry on connection errors
			if _, ok := err.(*net.OpError); ok {
				return true
			}

			return false
		}),
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

	response, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func methodIsCacheable(method string) bool {
	return method == http.MethodGet || method == http.MethodHead
}

func methodIsAction(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete
}
