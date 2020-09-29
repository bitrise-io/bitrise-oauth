package introspectors

import (
	"net/http"
	"time"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/labstack/echo"
	"gopkg.in/square/go-jose.v2"
)

// Token ...
type Token struct {
	validator *auth0.JWTValidator
}

// NewToken ...
func NewToken() Token {
	opts := auth0.JWKClientOptions{
		URI: "http://104.154.234.133/auth/realms/master/protocol/openid-connect/certs",
	}

	keyCacher := auth0.NewMemoryKeyCacher(3*time.Minute, 5)
	client := auth0.NewJWKClientWithCache(opts, nil, keyCacher)

	configuration := auth0.NewConfiguration(client, nil,
		"http://104.154.234.133/auth/realms/master", jose.RS256)

	return Token{
		validator: auth0.NewValidator(configuration, nil),
	}
}

// ValidateRequest ...
func (kti Token) ValidateRequest(r *http.Request) error {
	_, err := kti.validator.ValidateRequest(r)
	return err
}

// Middleware ...
func (kti Token) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := kti.ValidateRequest(r); err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MiddlewareFunc ...
func (kti Token) MiddlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if _, err := kti.validator.ValidateRequest(c.Request()); err != nil {
				return err
			}
			return next(c)
		}
	}
}

// HandlerFunc ...
func (kti Token) HandlerFunc(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := kti.ValidateRequest(r); err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		hf(w, r)
	}
}
