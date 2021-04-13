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
	asd, _ := jwt.ParseSigned(`eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ0dG80NS1qSzBFOURQa2huN2R2TVMzRUhWVXFsVUdoVERLcUJjdDlNNDRVIn0.eyJleHAiOjE2MTgyOTkwOTQsImlhdCI6MTYxODI5ODc5NCwianRpIjoiMTBjN2VkNWMtYWE5OC00NzRjLTgwOGItOWIxNDY4ZGM3Mjk4IiwiaXNzIjoiaHR0cHM6Ly9hdXRoLnNlcnZpY2VzLmJpdHJpc2UuZGV2L2F1dGgvcmVhbG1zL2FkZG9ucyIsImF1ZCI6ImJpdHJpc2UtYXBpIiwic3ViIjoiNmEyZDNiYTMtNzdlMS00MmJmLWE5NmItYzRkNmY2OTljOGVjIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiYml0cmlzZSIsInNlc3Npb25fc3RhdGUiOiIyYjJjMjMyOC1jYWEwLTQ0NjQtYTRkOS1mMzRjMzRlZTE2NTQiLCJhY3IiOiIxIiwic2NvcGUiOiJzZWxmOnJlYWQgYXBwOndyaXRlIGJ1aWxkOndyaXRlIGFwcDpyZWFkIGJ1aWxkOnJlYWQiLCJ1c2VyX2lkIjoiMzQ3Mjg1ZmYzMmVlZDE1OCJ9.Be1mLfucKpaxEqk7J0dOO9_qwtrUykKKZLATLAgtc0pz6gNqHSHyCQtKJyIxhF8mevEEgjrfsPlTkPqwervvQQfnzm6lDRXwQaojRUs-H7ZDnoMyD9IMClyussgDWSZsOlzLnp5mzBWto-eT74nZ4NSC53HHndNeRFMKOIEBVaWqVcr3KcfOP8YJm_Gk9VCQ5Woq9B2S9XaKx2zXFz90D6hEGFq5GCq39-g3rQSV_4duCZJxhZReZg8nsJ3Ju40CFhzXkvHPz7QR9X3kJsf4AQR9khewjYukBnWrMKjZwAUXpXmS7tnpzqztrJT0N2s5LKWadxIraV4gL5eQbfFBQQ`)
	m.On("ValidateRequest", mock.Anything).Return(asd, nil)
	return m
}

// GivenUnsuccessfulJWTValidation ...
func (m *JWTValidator) GivenUnsuccessfulJWTValidation(err error) *JWTValidator {
	m.On("ValidateRequest", mock.Anything).Return(&jwt.JSONWebToken{}, err)
	return m
}
