package authproviders

// ClientOption ...
type ClientOption func(c *ClientWithSecret)

// WithCustomTokenURL ...
func WithCustomTokenURL(tokenURL string) ClientOption {
	return func(c *ClientWithSecret) {
		c.tokenURL = tokenURL
	}
}
