package service

import "net/http"

// Introspector ...
type Introspector interface {
	Middleware(http.Handler) http.Handler
}
