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
		// mark error, custom error
	}

	return res, err
}
