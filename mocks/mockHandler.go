package mocks

import (
	"net/http"

	"github.com/c2fo/testify/mock"
)

// Handler ...
type Handler struct {
	mock.Mock
}

// ServeHTTP ...
func (m *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}
