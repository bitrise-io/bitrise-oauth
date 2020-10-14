package mocks

import (
	"net/http"

	"github.com/c2fo/testify/mock"
	"gopkg.in/square/go-jose.v2/jwt"
)

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
	m.On("ValidateRequest", mock.Anything).Return(&jwt.JSONWebToken{}, nil)
	return m
}

// GivenUnsuccessfulJTWValidation ...
func (m *JWTValidator) GivenUnsuccessfulJWTValidation(err error) *JWTValidator {
	m.On("ValidateRequest", mock.Anything).Return(&jwt.JSONWebToken{}, err)
	return m
}
