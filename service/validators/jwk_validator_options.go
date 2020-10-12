package validators

import "github.com/auth0-community/go-auth0"

// ValidatorOption ...
type ValidatorOption func(c *JWK)

// WithBaseURL ...
func WithBaseURL(url string) ValidatorOption {
	return func(c *JWK) {
		c.baseURL = url
	}
}

// WithRealm ...
func WithRealm(realm string) ValidatorOption {
	return func(c *JWK) {
		c.realm = realm
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
