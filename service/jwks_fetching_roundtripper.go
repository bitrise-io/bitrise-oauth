package service

import (
	"net/http"
)

// JWKSFetchingRoundTripper ...
type JWKSFetchingRoundTripper struct{}

// RoundTrip ...
func (roundTripper *JWKSFetchingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		err = &InternalError{Err: err}
	}

	return res, err
}
