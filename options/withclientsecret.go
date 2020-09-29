package options

import (
	"github.com/bitrise-io/bitriseoauth/settings"
)

// WithClientIDAndSecret ...
type WithClientIDAndSecret struct {
	ClientID, ClientSecret string
}

// Apply ...
func (opt WithClientIDAndSecret) Apply(settings *settings.ClientSettings) {
	// TODO: turn id and secret into a token
	// TODO: manage token storage
}
