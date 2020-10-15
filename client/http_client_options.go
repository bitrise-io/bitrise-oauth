package client

import (
	"context"
	"net/http"
)

// HTTPClientOption ...
type HTTPClientOption func(c *HTTPClientConfig)

// HTTPClientConfig ...
type HTTPClientConfig struct {
	context    context.Context
	baseClient *http.Client
}

// WithContext ...
func WithContext(ctx context.Context) HTTPClientOption {
	return func(c *HTTPClientConfig) {
		c.context = ctx
	}
}

// WithBaseClient ...
func WithBaseClient(bc *http.Client) HTTPClientOption {
	return func(c *HTTPClientConfig) {
		c.baseClient = bc
	}
}
