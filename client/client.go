package client

import "net/http"

// AuthProvider which creates an authenticated *http.Client ready to use with our authentication provider.
type AuthProvider interface {
	Client() *http.Client
}
