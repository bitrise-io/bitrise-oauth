package service

import "net/http"

// ExternalErrorHandler is a type alias for a function which recieves an error and returns an error
type ExternalErrorHandler func(error) error

// JWKSFetchingRoundTripper ...
type JWKSFetchingRoundTripper struct {
	ErrorHandler ExternalErrorHandler
	Base         http.RoundTripper
}

// RoundTrip ...
func (roundTripper *JWKSFetchingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := roundTripper.base().RoundTrip(req)

	if err != nil {
		roundTripper.ErrorHandler(err)
	}

	return res, err
}

func (roundTripper *JWKSFetchingRoundTripper) base() http.RoundTripper {
	if roundTripper.Base != nil {
		return roundTripper.Base
	}
	return http.DefaultTransport
}
