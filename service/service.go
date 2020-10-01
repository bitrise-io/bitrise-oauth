package service

import (
	"net/http"

	"github.com/labstack/echo"
)

// Validator gives multiple solution to validate the access token received in the request headers using Oauth2.0
type Validator interface {
	HandlerFunc(http.HandlerFunc) http.HandlerFunc
	Middleware(http.Handler) http.Handler
	MiddlewareFunc() echo.MiddlewareFunc
	ValidateRequest(r *http.Request) error
}
