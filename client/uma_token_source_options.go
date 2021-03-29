package client

import "github.com/bitrise-io/bitrise-oauth/config"

// UMATokenSourceOption represents a configuration option
type UMATokenSourceOption func(u *UMATokenSourceConfig)

// UMATokenSourceConfig represents the configuration of an
// UMATokenSource.
type UMATokenSourceConfig struct {
	audience *config.AudienceConfig
}

// WithAudienceConfig returns a function, which sets the audience
// to the provided audience configuration.
func WithAudienceConfig(c config.AudienceConfig) UMATokenSourceOption {
	return func(u *UMATokenSourceConfig) {
		u.audience = &c
	}
}
