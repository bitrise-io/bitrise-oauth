package client

import "net/http"

// AuthProvider ...
type AuthProvider interface {
	Client() *http.Client
}
