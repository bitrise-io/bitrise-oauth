package client

// Option ...
type Option func(c *WithSecret)

// WithTokenURL ...
func WithTokenURL(tokenURL string) Option {
	return func(c *WithSecret) {
		c.tokenURL = tokenURL
	}
}
