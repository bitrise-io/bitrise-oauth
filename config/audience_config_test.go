package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_All(t *testing.T) {
	config := NewAudienceConfig("aud1", "aud2")

	assert.Equal(t, []string{"aud2", "aud1"}, config.All())
}
