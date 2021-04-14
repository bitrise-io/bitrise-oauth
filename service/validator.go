package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/auth0-community/go-auth0"
	"github.com/bitrise-io/bitrise-oauth/config"
	"github.com/labstack/echo"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type jwtValidator interface {
	ValidateRequest(r *http.Request) (*jwt.JSONWebToken, error)
}

// Validator gives multiple solution to validate the access token received in the request headers using Oauth2.0
type Validator interface {
	HandlerFunc(http.HandlerFunc, ...HTTPMiddlewareOption) http.HandlerFunc
	Middleware(http.Handler, ...HTTPMiddlewareOption) http.Handler
	EchoMiddlewareFunc(...EchoMiddlewareOption) echo.MiddlewareFunc
	ValidateRequest(r *http.Request) error
	ValidateRequestAndReturnToken(r *http.Request) (TokenWithClaims, error)
}

// ValidatorConfig ...
type ValidatorConfig struct {
	jwtValidator       jwtValidator
	baseURL            string
	realm              string
	keyCacher          auth0.KeyCacher
	signatureAlgorithm jose.SignatureAlgorithm
	timeout            time.Duration
	audience           config.AudienceConfig
	issuer             string
	secretProvider     auth0.SecretProvider
}

// NewValidator returns the prepared JWK model. All input arguments are optional.
func NewValidator(audienceConfig config.AudienceConfig, opts ...ValidatorOption) Validator {
	serviceValidator := &ValidatorConfig{
		baseURL:            config.BaseURL,
		realm:              config.Realm,
		keyCacher:          auth0.NewMemoryKeyCacher(2*time.Hour, 5),
		signatureAlgorithm: jose.RS256,
		timeout:            30 * time.Second,
		audience:           audienceConfig,
	}

	for _, opt := range opts {
		opt(serviceValidator)
	}

	if len(serviceValidator.issuer) == 0 {
		serviceValidator.issuer = serviceValidator.realmURL()
	}

	if serviceValidator.secretProvider == nil {
		serviceValidator.secretProvider = createDefaultSecretProvider(serviceValidator)
	}

	if serviceValidator.jwtValidator == nil {
		serviceValidator.jwtValidator = createDefaultJWTValidator(serviceValidator)
	}

	return serviceValidator
}

func createDefaultSecretProvider(validatorConfig *ValidatorConfig) auth0.SecretProvider {
	secretProvderClientOptions := auth0.JWKClientOptions{
		URI:    validatorConfig.jwksURL(),
		Client: &http.Client{Timeout: validatorConfig.timeout},
	}

	return auth0.NewJWKClientWithCache(secretProvderClientOptions, nil, validatorConfig.keyCacher)
}

func createDefaultJWTValidator(validatorConfig *ValidatorConfig) jwtValidator {
	configuration := auth0.NewConfiguration(validatorConfig.secretProvider, []string{}, validatorConfig.issuer, validatorConfig.signatureAlgorithm)
	return auth0.NewValidator(configuration, nil)
}

func (sv ValidatorConfig) realmURL() string {
	return fmt.Sprintf("%s/auth/realms/%s", sv.baseURL, sv.realm)
}

func (sv ValidatorConfig) jwksURL() string {
	return fmt.Sprintf("%s/protocol/openid-connect/certs", sv.realmURL())
}

// ValidateRequest to validate if the request is authenticated and has active token.
func (sv ValidatorConfig) ValidateRequest(r *http.Request) error {
	token, err := sv.jwtValidator.ValidateRequest(r)

	key, err := sv.secretProvider.GetSecret(r)
	if err != nil {
		return err
	}

	tokenWithClaims := &tokenWithClaims{
		key:   key,
		token: token,
	}

	err = sv.validateAudiences(*tokenWithClaims, sv.audience.All())
	if err != nil {
		return err
	}

	return err
}

// ValidateRequestAndReturnToken ...
func (sv ValidatorConfig) ValidateRequestAndReturnToken(r *http.Request) (TokenWithClaims, error) {
	token, err := sv.jwtValidator.ValidateRequest(r)
	if err != nil {
		return nil, err
	}

	key, err := sv.secretProvider.GetSecret(r)
	if err != nil {
		return nil, err
	}

	tokenWithClaims := &tokenWithClaims{
		key:   key,
		token: token,
	}

	err = sv.validateAudiences(*tokenWithClaims, sv.audience.All())
	if err != nil {
		return nil, err
	}

	return tokenWithClaims, nil
}

// ValidateAudiences ...
func (sv ValidatorConfig) validateAudiences(tokenWithClaims tokenWithClaims, audiences []string) error {
	payload, err := tokenWithClaims.Payload()
	if err != nil {
		return err
	}

	var audiencesInToken []string

	switch v := payload["aud"].(type) {
	default:
		fmt.Printf("unexpected type %T", v)
	case string:
		audience, ok := payload["aud"].(string)
		if !ok {
			panic("type assertion failed")
		}
		audiencesInToken = []string{audience}
	case []interface{}:
		audienceInterface, ok := payload["aud"].([]interface{})
		if !ok {
			panic("type assertion failed")
		}
		audiencesInToken = make([]string, len(audienceInterface))
		for i, v := range audienceInterface {
			audience, ok := v.(string)
			if !ok {
				panic("type assertion failed")
			}
			audiencesInToken[i] = audience
		}
	}

	if len(audiencesInToken) > 0 {
		found := false
		for _, aud := range audiencesInToken {
			if !found && sv.audience.Contains(aud) {
				found = true
			}
		}

		if !found {
			return jwt.ErrInvalidAudience
		}
	}

	return nil
}

// Middleware used as http package's middleware, in http.Handle.
// Calls out to ValidateRequest and returns http.Status Unauthorized with body: invalid token if the token is not active.
func (sv ValidatorConfig) Middleware(next http.Handler, opts ...HTTPMiddlewareOption) http.Handler {
	handlerConfig := &HTTPMiddlewareConfig{
		errorWriter: defaultHTTPErrorWriter,
	}

	for _, opt := range opts {
		opt(handlerConfig)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := sv.ValidateRequestAndReturnToken(r)
		if err != nil {
			handlerConfig.errorWriter(w, r, err)
			return
		}
		if handlerConfig.tokenHandler != nil {
			handlerConfig.tokenHandler(w, r, token)
		}
		next.ServeHTTP(w, r)
	})
}

// EchoMiddlewareFunc can be used with echo.Use.
// Calls out to ValidateRequest and returns an error for echo.
func (sv ValidatorConfig) EchoMiddlewareFunc(opts ...EchoMiddlewareOption) echo.MiddlewareFunc {
	handlerConfig := &EchoMiddlewareConfig{
		errorWriter: defaultEchoErrorWriter,
	}

	for _, opt := range opts {
		opt(handlerConfig)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := sv.ValidateRequest(c.Request()); err != nil {
				return handlerConfig.errorWriter(c, err)
			}
			return next(c)
		}
	}
}

// HandlerFunc used with http.HandleFunc.
// Calls out to ValidateRequest and returns http.Status Unauthorized with body: invalid token if the token is not active.
func (sv ValidatorConfig) HandlerFunc(hf http.HandlerFunc, opts ...HTTPMiddlewareOption) http.HandlerFunc {
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
