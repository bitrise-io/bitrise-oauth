package mocks

import (
	"github.com/c2fo/testify/mock"
)

// InternalErrorHandler ...
type InternalErrorHandler struct {
	mock.Mock
}

// HandlerFunction ...
func (m *InternalErrorHandler) HandlerFunction(err error) {
	m.Called(err)
}
