package mocks

import (
	"net/http"

	"github.com/c2fo/testify/mock"
	"gopkg.in/square/go-jose.v2/jwt"
)

// token can be edited on jwt.io (use "secret" as the secret)
var rawMockToken = `eyJhbGciOiJIUzI1NiJ9.eyJTY29wZXMiOlsiZm9vIiwiYmFyIl0sImlzcyI6Imlzc3VlciIsInN1YiI6InN1YmplY3QiLCJhdWQiOiJ0ZXN0X2F1ZGllbmNlIn0.jWFz6fxqVOWZOUEj50_FjKIvZZRjjAOxk5YgpXg0aLI`
var mockToken, _ = jwt.ParseSigned(rawMockToken)

// JWTValidator ...
type JWTValidator struct {
	mock.Mock
}

// ValidateRequest ...
func (m *JWTValidator) ValidateRequest(r *http.Request) (*jwt.JSONWebToken, error) {
	args := m.Called(r)
	return args.Get(0).(*jwt.JSONWebToken), args.Error(1)
}

// GivenSuccessfulJWTValidation ...
func (m *JWTValidator) GivenSuccessfulJWTValidation() *JWTValidator {
	m.On("ValidateRequest", mock.Anything).Return(mockToken, nil)
	return m
}

// GivenUnsuccessfulJWTValidation ...
func (m *JWTValidator) GivenUnsuccessfulJWTValidation(err error) *JWTValidator {
	m.On("ValidateRequest", mock.Anything).Return(&jwt.JSONWebToken{}, err)
	return m
}
