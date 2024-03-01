package requestbuilder

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cybre/home-inventory/internal/cache"
	"github.com/cybre/home-inventory/internal/logging"
	"github.com/pquerna/cachecontrol"
	"github.com/pquerna/cachecontrol/cacheobject"
)

type AutoCacheTransport struct {
	Transport http.RoundTripper
	Cache     *cache.Cache[string]
	CacheKey  string
}

func NewAutoCacheTransport(transport http.RoundTripper, cache *cache.Cache[string], cacheKey string) AutoCacheTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return AutoCacheTransport{
		Transport: transport,
		Cache:     cache,
		CacheKey:  cacheKey,
	}
}

func (t AutoCacheTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	logger := logging.FromContext(req.Context())

	reqDir, err := cacheobject.ParseRequestCacheControl(req.Header.Get("Cache-Control"))
	if err != nil {
		logger.Warn("failed to parse request cache control", slog.Any("error", err))
	}

	if reqDir == nil || !reqDir.NoCache {
		if value, err := t.Cache.Get(req.Context(), t.CacheKey); err == nil {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(value)),
			}, nil
		}
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	reasons, expires, _ := cachecontrol.CachableResponse(req, resp, cachecontrol.Options{
		PrivateCache: false,
	})

	if len(reasons) > 0 {
		return resp, nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("failed to read response body: %w", err)
	}

	if expires.IsZero() {
		if err := t.Cache.Set(req.Context(), t.CacheKey, string(body)); err != nil {
			logger.Warn("failed to set response cache", slog.Any("error", err))
		}
	} else if err := t.Cache.SetWithExpiration(req.Context(), t.CacheKey, string(body), expires.Sub(time.Now())); err != nil {
		logger.Warn("failed to set response cache with expiration", slog.Any("error", err))
	}

	resp.Body = io.NopCloser(bytes.NewReader(body))

	return resp, nil
}
