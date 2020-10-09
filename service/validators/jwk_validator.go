package validators

import (
	"net/http"
	"time"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/bitrise-io/bitrise-oauth/config"
	"github.com/bitrise-io/bitrise-oauth/service"
	"github.com/labstack/echo"
	"gopkg.in/square/go-jose.v2"
)

// JWK ...
type JWK struct {
	validator *auth0.JWTValidator
	baseURL   string
	realm     string
	keyCacher auth0.KeyCacher
	jwksURL   string
	realmURL  string
}

// NewJWK returns the prepared JWK model. All input arguments are optional.
//
// Argument defaults when nil:
//  	baseURL: http://104.154.234.133
//  	realm: master
//  	keyCacher: auth0 MemoryKeyCacher with 3 minutes TTL and size 5
func NewJWK(opts ...ValidatorOption) service.Validator {
	serviceValidator := &JWK{
		baseURL:   config.BaseURL,
		realm:     config.Realm,
		keyCacher: auth0.NewMemoryKeyCacher(3*time.Minute, 5),
		jwksURL:   config.JWKSURL,
		realmURL:  config.RealmURL,
	}

	for _, opt := range opts {
		opt(serviceValidator)
	}

	clientOpts := auth0.JWKClientOptions{
		URI: serviceValidator.jwksURL,
	}

	client := auth0.NewJWKClientWithCache(clientOpts, nil, serviceValidator.keyCacher)

	configuration := auth0.NewConfiguration(client, nil,
		serviceValidator.realmURL, jose.RS256)

	serviceValidator.validator = auth0.NewValidator(configuration, nil)

	return serviceValidator
}

// ValidateRequest to validate if the request is authenticated and has active token.
func (kti JWK) ValidateRequest(r *http.Request) error {
	_, err := kti.validator.ValidateRequest(r)
	return err
}

// Middleware used as http package's middleware, in http.Handle.
// Calls out to ValidateRequest and returns http.StatusUnauthorized with body: invalid token if the token is not active.
func (kti JWK) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := kti.ValidateRequest(r); err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MiddlewareFunc can be used with echo.Use.
// Calls out to ValidateRequest and returns an error for echo.
func (kti JWK) MiddlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if _, err := kti.validator.ValidateRequest(c.Request()); err != nil {
				return err
			}
			return next(c)
		}
	}
}

// HandlerFunc used with http.HandleFunc.
// Calls out to ValidateRequest and returns http.StatusUnauthorized with body: invalid token if the token is not active.
func (kti JWK) HandlerFunc(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := kti.ValidateRequest(r); err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		hf(w, r)
	}
}
