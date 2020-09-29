package authprovider

import "net/http"

// Handler ...
type Handler interface {
	ClientWithIDAndSecret(id, secret string) *http.Client
}
