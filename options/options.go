package options

import "github.com/bitrise-io/bitriseoauth/settings"

// ClientOption ...
type ClientOption interface {
	Apply(settings *settings.ClientSettings)
}
