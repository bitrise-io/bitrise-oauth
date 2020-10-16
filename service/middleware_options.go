package service

import (
	"net/http"

	"github.com/labstack/echo"
)

var defaultHTTPErrorWriter = func(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
}

// HTTPMiddlewareOption ...
type HTTPMiddlewareOption func(c *HTTPMiddlewareConfig)

// HTTPMiddlewareConfig ...
type HTTPMiddlewareConfig struct {
	errorWriter func(w http.ResponseWriter, r *http.Request, err error)
}

// WithHTTPErrorWriter ...
func WithHTTPErrorWriter(errorWriter func(w http.ResponseWriter, r *http.Request, err error)) HTTPMiddlewareOption {
	return func(c *HTTPMiddlewareConfig) {
		c.errorWriter = errorWriter
	}
}

var defaultEchoErrorWriter = func(c echo.Context, err error) error {
	return err
}

// EchoMiddlewareOption ...
type EchoMiddlewareOption func(c *EchoMiddlewareConfig)

// EchoMiddlewareConfig ...
type EchoMiddlewareConfig struct {
	errorWriter func(echo.Context, error) error
}

// WithContextErrorWriter ...
func WithContextErrorWriter(errorWriter func(echo.Context, error) error) EchoMiddlewareOption {
	return func(c *EchoMiddlewareConfig) {
		c.errorWriter = errorWriter
	}
}
