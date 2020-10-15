package validators

import (
	"net/http"

	"github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
)

// ValidatorOption ...
type ValidatorOption func(c *JWK)

// WithBaseURL ...
func WithBaseURL(url string) ValidatorOption {
	return func(c *JWK) {
		c.baseURL = url
	}
}

// WithSignatureAlgorithm ...
func WithSignatureAlgorithm(sa jose.SignatureAlgorithm) ValidatorOption {
	return func(c *JWK) {
		c.signatureAlgorithm = sa
	}
}

// WithRealm ...
func WithRealm(realm string) ValidatorOption {
	return func(c *JWK) {
		c.realm = realm
	}
}

// WithErrorWriter ...
func WithErrorWriter(errorWriter func(http.ResponseWriter)) ValidatorOption {
	return func(c *JWK) {
		c.errorWriter = errorWriter
	}
}

// WithKeyCacher ...
func WithKeyCacher(kc auth0.KeyCacher) ValidatorOption {
	return func(c *JWK) {
		c.keyCacher = kc
	}
}

// WithRealmURL ...
func WithRealmURL(realmURL string) ValidatorOption {
	return func(c *JWK) {
		c.realmURL = realmURL
	}
}

// WithJWKSURL ...
func WithJWKSURL(jwksURL string) ValidatorOption {
	return func(c *JWK) {
		c.jwksURL = jwksURL
	}
}

// WithValidator ...
func WithValidator(validator JWTValidator) ValidatorOption {
	return func(c *JWK) {
		c.validator = validator
	}
}
