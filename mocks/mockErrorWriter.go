package mocks

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

// ErrorWriter ...
type ErrorWriter struct {
	mock.Mock
}

// ErrorHandler ...
func (m *ErrorWriter) ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	m.Called(w, r, err)
}

// EchoHandlerFunc ...
func (m *ErrorWriter) EchoHandlerFunc(c echo.Context, err error) error {
	args := m.Called(c, err)
	return args.Error(0)
}
