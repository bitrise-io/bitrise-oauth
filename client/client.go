package client

import (
	"net/http"

	"github.com/bitrise-io/bitriseoauth/authprovider"
	"github.com/bitrise-io/bitriseoauth/options"
)

var provider authprovider.Hydra

// New ...
func New(options ...options.ClientOption) *http.Client {
	settings := provider.ClientSettings()

	for _, option := range options {
		option.Apply(&settings)
	}

	return http.DefaultClient
}
