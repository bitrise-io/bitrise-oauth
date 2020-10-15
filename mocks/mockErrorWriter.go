package mocks

import (
	"net/http"

	"github.com/c2fo/testify/mock"
)

// ErrorWriter ...
type ErrorWriter struct {
	mock.Mock
}

// ErrorHandler ...
func (m *ErrorWriter) ErrorHandler(w http.ResponseWriter) {
	m.Called(w)
}
