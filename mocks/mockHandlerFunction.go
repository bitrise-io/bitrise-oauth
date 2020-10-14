package mocks

import (
	"net/http"

	"github.com/c2fo/testify/mock"
)

// HandlerFunction ...
type HandlerFunction struct {
	mock.Mock
}

// Handler ...
func (m *HandlerFunction) Handler(w http.ResponseWriter, r *http.Request) {
	m.Called()
}
