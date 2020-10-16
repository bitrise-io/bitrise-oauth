package mocks

import (
	"github.com/c2fo/testify/mock"
	"github.com/labstack/echo"
)

// MiddlewareHandlerFunction ...
type MiddlewareHandlerFunction struct {
	mock.Mock
}

// HandlerFunction ...
func (m *MiddlewareHandlerFunction) HandlerFunction(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

// GivenSuccess ...
func (m *MiddlewareHandlerFunction) GivenSuccess() *MiddlewareHandlerFunction {
	m.On("HandlerFunction", mock.Anything).Return(nil)
	return m
}

// GivenError ...
func (m *MiddlewareHandlerFunction) GivenError(err error) *MiddlewareHandlerFunction {
	m.On("HandlerFunction", mock.Anything).Return(err)
	return m
}
