package client

import (
	"net/http"

	"golang.org/x/oauth2"
)

type invalidTokenRefresherTransport struct {
	base     *oauth2.Transport
	tokenSrc oauth2.TokenSource
}

func (t *invalidTokenRefresherTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		if resettableTokenSource, ok := t.tokenSrc.(*resettableTokenSource); ok {
			resettableTokenSource.Reset()
			return t.base.RoundTrip(req)
		}
	}

	return resp, err
}
