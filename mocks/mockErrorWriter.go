package mocks

import (
	"net/http"

	"github.com/c2fo/testify/mock"
	"github.com/labstack/echo"
)

// ErrorWriter ...
type ErrorWriter struct {
	mock.Mock
}

// ErrorHandler ...
func (m *ErrorWriter) ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	m.Called(w, nil, err)
}

// EchoHandlerFunc ...
func (m *ErrorWriter) EchoHandlerFunc(c echo.Context, err error) error {
	args := m.Called(c, err)
	return args.Error(0)
}
