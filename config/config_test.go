package config_test

import (
	"testing"

	"github.com/bitrise-io/bitrise-oauth/config"
	"github.com/c2fo/testify/assert"
)

func Test_DefaultConfigValues(t *testing.T) {
	defaultBaseURL := "https://auth.services.bitrise.io"
	assert.Equal(t, defaultBaseURL, config.BaseURL)

	defaultRealm := "bitrise-services"
	assert.Equal(t, defaultRealm, config.Realm)
}
