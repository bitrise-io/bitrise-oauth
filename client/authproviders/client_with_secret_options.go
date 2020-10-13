package authproviders

// ClientOption ...
type ClientOption func(c *ClientWithSecret)

// WithTokenURL ...
func WithTokenURL(tokenURL string) ClientOption {
	return func(c *ClientWithSecret) {
		c.tokenURL = tokenURL
	}
}
