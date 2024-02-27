package requestbuilder_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/cybre/home-inventory/internal/requestbuilder"
	"github.com/stretchr/testify/assert"
)

func Test_RequestBuilder_Build(t *testing.T) {
	// Create a new context
	ctx := context.WithValue(context.Background(), "correlation_id", "123")

	// Create a slice of multiple values for a query parameter
	multipleValues := []string{"value1", "value2"}

	// Create a request body
	body := struct {
		Key string `json:"key"`
	}{
		Key: "value",
	}

	// Create a new request builder
	req, err := requestbuilder.New(http.MethodGet, "http://example.com/:pathParam1/:pathParam2").
		// Add headers to the request
		WithHeader("Cache-Control", "max-age=0").
		WithHeader("Connection", "keep-alive").
		// Add a query parameter to the request
		WithQueryParam("single", "value").
		// Add query parameter with multiple values to the request
		WithQueryParam("multiple", multipleValues...).
		// Add a path parameters to the request
		WithPathParam("pathParam1", "pathValue1").
		WithPathParam("pathParam2", "pathValue2").
		// Add a request body to the request
		WithBody(body).
		Build(ctx)

	// Validate the request
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.True(t, strings.HasPrefix(req.URL.String(), "http://example.com/pathValue1/pathValue2"))
	assert.Equal(t, "max-age=0", req.Header.Get("Cache-Control"))
	assert.Equal(t, "keep-alive", req.Header.Get("Connection"))
	assert.Equal(t, "value", req.URL.Query().Get("single"))
	assert.ElementsMatch(t, multipleValues, req.URL.Query()["multiple"])
	if assert.NotNil(t, req.Body) {
		defer req.Body.Close()
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		var reqBody struct {
			Key string `json:"key"`
		}
		assert.NoError(t, json.NewDecoder(req.Body).Decode(&reqBody))
		assert.Equal(t, body, reqBody)
	}
	assert.Equal(t, ctx, req.Context())
	assert.Equal(t, "123", req.Header.Get("X-Correlation-ID"))
}
