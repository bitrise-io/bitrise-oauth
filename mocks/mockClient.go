package mocks

import "github.com/stretchr/testify/mock"

// Client ...
type Client struct {
	mock.Mock
}

// Test ...
func (m *Client) Test(accessToken string) {
	m.Called(accessToken)
}
