package mocks

import (
	"net/http"

	"github.com/c2fo/testify/mock"
	"gopkg.in/square/go-jose.v2/jwt"
)

// HandlerFunction ...
type HandlerFunction struct {
	mock.Mock
}

// Handler ...
func (m *HandlerFunction) Handler(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// GivenNoReturnValue ...
func (m *HandlerFunction) GivenNoReturnValue() *HandlerFunction {
	m.On("Handler", mock.Anything, mock.Anything).Return(&jwt.JSONWebToken{}, nil)
	return m
}
