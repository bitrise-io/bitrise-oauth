package mocks

import "github.com/c2fo/testify/mock"

// Client ...
type Client struct {
	mock.Mock
}

// Test ...
func (m *Client) Test(accessToken string) {
	m.Called(accessToken)
}
