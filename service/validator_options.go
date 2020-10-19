package service

import (
	"github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
)

// ValidatorOption ...
type ValidatorOption func(c *ValidatorConfig)

// WithBaseURL ...
func WithBaseURL(url string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.baseURL = url
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

// WithRealmURL ...
func WithRealmURL(realmURL string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.realmURL = realmURL
	}
}

// WithJWKSURL ...
func WithJWKSURL(jwksURL string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.jwksURL = jwksURL
	}
}

// WithValidator ...
func WithValidator(validator JWTValidator) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.jwtValidator = validator
	}
}
