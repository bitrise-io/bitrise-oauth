package mocks

import "github.com/c2fo/testify/mock"

// AuthService ...
type AuthService struct {
	mock.Mock
}

// Token simulates the /token endpoint called
func (m *AuthService) Token(grantType string) {
	m.Called(grantType)
}

// Certs simulates the /certs endpoint called
func (m *AuthService) Certs() {
	m.Called()
}
