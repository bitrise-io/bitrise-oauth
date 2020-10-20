package client

// Option ...
type Option func(c *WithSecret)

// WithBaseURL for example: https://auth.services.bitrise.io
func WithBaseURL(baseURL string) Option {
	return func(c *WithSecret) {
		c.baseURL = baseURL
	}
}

// WithRealm for example: master
func WithRealm(realm string) Option {
	return func(c *WithSecret) {
		c.realm = realm
	}
}
