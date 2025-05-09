package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// Handler ...
type Handler struct {
	mock.Mock
}

// ServeHTTP ...
func (m *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}
