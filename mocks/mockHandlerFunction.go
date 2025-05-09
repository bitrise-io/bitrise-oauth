package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// HandlerFunction ...
type HandlerFunction struct {
	mock.Mock
}

// Handler ...
func (m *HandlerFunction) Handler(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}
