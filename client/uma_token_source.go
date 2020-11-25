package client

import (
	"golang.org/x/oauth2"
)

// UMATokenSource returns a token source that returns a new token each time the Token()
// method is called.
type UMATokenSource interface {
	Token() (*oauth2.Token, error)
}

type umaTokenSource struct{}

// NewUMATokenSource returns a new UMA token source.
func NewUMATokenSource() UMATokenSource {
	return umaTokenSource{}
}

// Token returns a new UMA token upon each invocation.
func (tokenSource umaTokenSource) Token() (*oauth2.Token, error) {

	return nil, nil
}
