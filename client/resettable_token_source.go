package client

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type resettableTokenSource struct {
	creds clientcredentials.Config
	src   oauth2.TokenSource
	ctx   context.Context
}

func (rts *resettableTokenSource) Token() (*oauth2.Token, error) {
	return rts.src.Token()
}

func (rts *resettableTokenSource) Reset() {
	rts.src = rts.creds.TokenSource(rts.ctx)
}
