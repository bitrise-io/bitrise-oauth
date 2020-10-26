package service

import (
	"strings"

	"github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
)

// ValidatorOption ...
type ValidatorOption func(c *ValidatorConfig)

var defaultInternalErrorHandler = func(err error) {}

// WithBaseURL ...
func WithBaseURL(url string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.baseURL = strings.TrimSuffix(url, "/")
	}
}

// WithSignatureAlgorithm ...
func WithSignatureAlgorithm(sa jose.SignatureAlgorithm) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.signatureAlgorithm = sa
	}
}

// WithRealm ...
func WithRealm(realm string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.realm = realm
	}
}

// WithKeyCacher ...
func WithKeyCacher(kc auth0.KeyCacher) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.keyCacher = kc
	}
}

// WithValidator ...
func WithValidator(validator JWTValidator) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.jwtValidator = validator
	}
}

// WithInternalErrorHandler ...
func WithInternalErrorHandler(handler InternalErrorHandler) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.internalErrorHandler = handler
	}
}
