package requestbuilder

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/cybre/home-inventory/internal/cache"
	"github.com/cybre/home-inventory/internal/logging"
)

type SetCacheFn func() (string, any)

type SetCacheTransport struct {
	Transport http.RoundTripper
	Cache     *cache.Cache[string]
	SetFns    []SetCacheFn
}

func NewSetCacheTransport(transport http.RoundTripper, cache *cache.Cache[string], setFns ...SetCacheFn) SetCacheTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return SetCacheTransport{
		Transport: transport,
		Cache:     cache,
		SetFns:    setFns,
	}
}

func (t SetCacheTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !methodIsAction(req.Method) {
		return t.Transport.RoundTrip(req)
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, nil
	}

	logger := logging.FromContext(req.Context())

	for _, fn := range t.SetFns {
		key, value := fn()

		if value == nil {
			if err := t.Cache.Delete(req.Context(), key); err != nil {
				logger.Error("failed to delete cache via setCacheFn", slog.Any("error", err))
			}

			continue
		}

		stringValue, err := json.Marshal(value)
		if err != nil {
			logger.Error("failed to marshal SetCacheFn return value", slog.Any("error", err))
			continue
		}

		if err := t.Cache.Set(req.Context(), key, string(stringValue)); err != nil {
			logger.Error("failed to set SetCacheFn cache", slog.Any("error", err))
		}
	}

	return resp, nil
}
