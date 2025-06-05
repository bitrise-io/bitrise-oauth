package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_All(t *testing.T) {
	config := NewAudienceConfig("aud1", "aud2")

	assert.Equal(t, []string{"aud2", "aud1"}, config.All())
}

func Test_NewAudienceConfigFromAudiences(t *testing.T) {
	audiences := []string{"aud1", "aud2"}
	config := NewAudienceConfigFromAudiences(audiences)

	assert.Equal(t, audiences, config.All())
}
