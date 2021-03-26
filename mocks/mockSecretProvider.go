package mocks

import (
	"net/http"

	"github.com/auth0-community/go-auth0"
	"github.com/c2fo/testify/mock"
)

var _ auth0.SecretProvider = &MockSecretProvider{}

// MockSecretProvider ...
type MockSecretProvider struct {
	mock.Mock
}

// ValidateRequest ...
func (m *MockSecretProvider) GetSecret(r *http.Request) (interface{}, error) {
	args := m.Called(r)
	return args.Get(0), args.Error(1)
}
