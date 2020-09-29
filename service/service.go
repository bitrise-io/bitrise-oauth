package service

import (
	"net/http"

	"github.com/labstack/echo"
)

// Introspector ...
type Introspector interface {
	HandlerFunc(http.HandlerFunc) http.HandlerFunc
	Middleware(http.Handler) http.Handler
	MiddlewareFunc() echo.MiddlewareFunc
	ValidateRequest(r *http.Request) error
}
