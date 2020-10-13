package authproviders

// ClientOption ...
type ClientOption func(c *ClientWithSecret)

// WithTokenURL ...
func WithTokenURL(tokenURL string) ClientOption {
	return func(c *ClientWithSecret) {
		c.tokenURL = tokenURL
	}
}

// WithCondition ...
func WithCondition(condition func() bool) ClientOption {
	return func(c *ClientWithSecret) {
		c.condition = condition
	}
}
