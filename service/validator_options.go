package service

import (
	"strings"
	"time"

	"github.com/bitrise-io/go-auth0"
	"gopkg.in/go-jose/go-jose.v2"
)

// ValidatorOption ...
type ValidatorOption func(c *ValidatorConfig)

// WithBaseURL ...
func WithBaseURL(url string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.baseURL = strings.TrimSuffix(url, "/")
	}
}

// WithTimeout ...
func WithTimeout(timeout time.Duration) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.timeout = timeout
	}
}

// WithSignatureAlgorithm ...
func WithSignatureAlgorithm(sa jose.SignatureAlgorithm) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.signatureAlgorithm = sa
	}
}

// WithIssuer ...
func WithIssuer(issuer string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.issuer = issuer
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

func withValidator(validator jwtValidator) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.jwtValidator = validator
	}
}

func withIssuer(issuer string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.issuer = issuer
	}
}

func withSecretProvider(secretProvider auth0.SecretProvider) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.secretProvider = secretProvider
	}
}
