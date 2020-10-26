package service

import "net/http"

// InternalErrorHandler is a type alias for a function that receives an error and returns an error.
// The caller can be notified about the internal errors of the package.
type InternalErrorHandler func(error)

// JWKSFetchingRoundTripper ...
type JWKSFetchingRoundTripper struct {
	ErrorHandler InternalErrorHandler
	Base         http.RoundTripper
}

// RoundTrip ...
func (roundTripper *JWKSFetchingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		roundTripper.ErrorHandler(err)
	}

	return res, err
}
