package validators

import (
	"github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
)

// ValidatorOption ...
type ValidatorOption func(c *JWK)

// WithCustomBaseURL ...
func WithCustomBaseURL(url string) ValidatorOption {
	return func(c *JWK) {
		c.baseURL = url
	}
}

// WithCustomSignatureAlgorithm ...
func WithCustomSignatureAlgorithm(sa jose.SignatureAlgorithm) ValidatorOption {
	return func(c *JWK) {
		c.signatureAlgorithm = sa
	}
}

// WithCustomRealm ...
func WithCustomRealm(realm string) ValidatorOption {
	return func(c *JWK) {
		c.realm = realm
	}
}

// WithCustomKeyCacher ...
func WithCustomKeyCacher(kc auth0.KeyCacher) ValidatorOption {
	return func(c *JWK) {
		c.keyCacher = kc
	}
}

// WithCustomRealmURL ...
func WithCustomRealmURL(realmURL string) ValidatorOption {
	return func(c *JWK) {
		c.realmURL = realmURL
	}
}

// WithCustomJWKSURL ...
func WithCustomJWKSURL(jwksURL string) ValidatorOption {
	return func(c *JWK) {
		c.jwksURL = jwksURL
	}
}
