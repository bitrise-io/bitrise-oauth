package client

import (
	"strings"
)

// Option ...
type Option func(c *WithSecret)

// WithBaseURL for example: https://auth.services.bitrise.io
func WithBaseURL(baseURL string) Option {
	return func(c *WithSecret) {
		c.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// WithRealm for example: master
func WithRealm(realm string) Option {
	return func(c *WithSecret) {
		c.realm = realm
	}
}

// ScopeOption ...
type ScopeOption func(c *WithSecret)

// WithScope ...
func WithScope(aud string, auds ...string) ScopeOption {
	return func(c *WithSecret) {
		c.scopes = append(auds, aud)
	}
}
