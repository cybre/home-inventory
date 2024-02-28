package requestbuilder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/rehttp"
	"github.com/cybre/home-inventory/internal/cache"
	"github.com/cybre/home-inventory/internal/logging"
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
	cache       *cache.Cache[string]
	cacheKey    string
	setCacheFn  func() any
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
		cache:       nil,
		cacheKey:    "",
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

func (r *RequestBuilder) WithCache(cache *cache.Cache[string], key string) *RequestBuilder {
	if !methodIsCacheable(r.method) {
		panic("WithCache can only be used with GET and HEAD requests")
	}

	r.cache = cache
	r.cacheKey = key

	return r
}

func (r *RequestBuilder) WithSetCacheFn(cache *cache.Cache[string], key string, setCacheFn func() any) *RequestBuilder {
	if !methodIsAction(r.method) {
		panic("WithSetCacheFn can only be used with POST, PUT, PATCH, and DELETE requests")
	}

	r.cache = cache
	r.cacheKey = key
	r.setCacheFn = setCacheFn

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

func (r *RequestBuilder) getFromCache(ctx context.Context) (*http.Response, bool) {
	if r.cache == nil || !methodIsCacheable(r.method) {
		return nil, false
	}

	value, err := r.cache.Get(ctx, r.cacheKey)
	if err != nil {
		return nil, false
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(value)),
	}, true
}

func (r *RequestBuilder) saveToCache(ctx context.Context, response *http.Response) error {
	if r.cache == nil || (response.StatusCode < 200 || response.StatusCode >= 300) {
		return nil
	}

	if r.setCacheFn != nil && !methodIsCacheable(r.method) {
		value := r.setCacheFn()
		if value == nil {
			if err := r.cache.Delete(ctx, r.cacheKey); err != nil {
				return fmt.Errorf("failed to delete cache via setCacheFn: %w", err)
			}
		}

		stringValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal setCacheFn return value: %w", err)
		}

		if err := r.cache.Set(ctx, r.cacheKey, string(stringValue)); err != nil {
			return fmt.Errorf("failed to set cacheFn cache: %w", err)
		}

		return nil
	}

	if methodIsCacheable(r.method) {
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if err := r.cache.Set(ctx, r.cacheKey, string(body)); err != nil {
			return fmt.Errorf("failed to set response cache: %w", err)
		}

		response.Body = io.NopCloser(bytes.NewReader(body))
	}

	return nil
}

func (r *RequestBuilder) Do(ctx context.Context) (*http.Response, error) {
	if resp, hit := r.getFromCache(ctx); hit {
		return resp, nil
	}

	req, err := r.Build(ctx)
	if err != nil {
		return nil, err
	}

	response, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err := r.saveToCache(ctx, response); err != nil {
		logging.FromContext(ctx).Error("failed to save response body to cache", slog.Any("error", err))
	}

	return response, nil
}

func methodIsCacheable(method string) bool {
	return method == http.MethodGet || method == http.MethodHead
}

func methodIsAction(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete
}
