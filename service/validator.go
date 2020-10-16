package service

import (
	"net/http"
	"time"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/bitrise-io/bitrise-oauth/config"
	"github.com/labstack/echo"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// JWTValidator ...
type JWTValidator interface {
	ValidateRequest(r *http.Request) (*jwt.JSONWebToken, error)
}

// ValidatorIntf gives multiple solution to validate the access token received in the request headers using Oauth2.0
type ValidatorIntf interface {
	HandlerFunc(http.HandlerFunc, ...HTTPMiddlewareOption) http.HandlerFunc
	Middleware(http.Handler, ...HTTPMiddlewareOption) http.Handler
	MiddlewareFunc(...EchoMiddlewareOption) echo.MiddlewareFunc
	ValidateRequest(r *http.Request) error
}

// Validator ...
type Validator struct {
	validator          JWTValidator
	baseURL            string
	realm              string
	keyCacher          auth0.KeyCacher
	jwksURL            string
	realmURL           string
	signatureAlgorithm jose.SignatureAlgorithm
}

// NewValidator returns the prepared JWK model. All input arguments are optional.
//
// Argument defaults when nil:
//  	baseURL: http://104.154.234.133
//  	realm: master
//  	keyCacher: auth0 MemoryKeyCacher with 3 minutes TTL and size 5
func NewValidator(opts ...ValidatorOption) ValidatorIntf {
	serviceValidator := &Validator{
		baseURL:            config.BaseURL,
		realm:              config.Realm,
		keyCacher:          auth0.NewMemoryKeyCacher(3*time.Minute, 5),
		jwksURL:            config.JWKSURL,
		realmURL:           config.RealmURL,
		signatureAlgorithm: jose.RS256,
	}

	for _, opt := range opts {
		opt(serviceValidator)
	}

	if serviceValidator.validator == nil {
		clientOpts := auth0.JWKClientOptions{
			URI: serviceValidator.jwksURL,
		}

		client := auth0.NewJWKClientWithCache(clientOpts, nil, serviceValidator.keyCacher)

		configuration := auth0.NewConfiguration(client, nil,
			serviceValidator.realmURL, serviceValidator.signatureAlgorithm)

		serviceValidator.validator = auth0.NewValidator(configuration, nil)
	}

	return serviceValidator
}

// ValidateRequest to validate if the request is authenticated and has active token.
func (sv Validator) ValidateRequest(r *http.Request) error {
	_, err := sv.validator.ValidateRequest(r)
	return err
}

// Middleware used as http package's middleware, in http.Handle.
// Calls out to ValidateRequest and returns http.StatusUnauthorized with body: invalid token if the token is not active.
func (sv Validator) Middleware(next http.Handler, opts ...HTTPMiddlewareOption) http.Handler {
	handlerConfig := &HTTPMiddlewareConfig{
		errorWriter: defaultHTTPErrorWriter,
	}

	for _, opt := range opts {
		opt(handlerConfig)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := sv.ValidateRequest(r); err != nil {
			handlerConfig.errorWriter(w, r, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MiddlewareFunc can be used with echo.Use.
// Calls out to ValidateRequest and returns an error for echo.
func (sv Validator) MiddlewareFunc(opts ...EchoMiddlewareOption) echo.MiddlewareFunc {
	handlerConfig := &EchoMiddlewareConfig{
		errorWriter: defaultEchoErrorWriter,
	}

	for _, opt := range opts {
		opt(handlerConfig)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if _, err := sv.validator.ValidateRequest(c.Request()); err != nil {
				return handlerConfig.errorWriter(c, err)
			}
			return next(c)
		}
	}
}

// HandlerFunc used with http.HandleFunc.
// Calls out to ValidateRequest and returns http.StatusUnauthorized with body: invalid token if the token is not active.
func (sv Validator) HandlerFunc(hf http.HandlerFunc, opts ...HTTPMiddlewareOption) http.HandlerFunc {
	handlerConfig := &HTTPMiddlewareConfig{
		errorWriter: defaultHTTPErrorWriter,
	}

	for _, opt := range opts {
		opt(handlerConfig)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := sv.ValidateRequest(r); err != nil {
			handlerConfig.errorWriter(w, r, err)
			return
		}
		hf(w, r)
	}
}
