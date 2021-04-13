package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_All(t *testing.T) {
	config := NewAudienceConfig("aud1", "aud2")

	assert.Equal(t, []string{"aud2", "aud1"}, config.All())
}

func Test_Contains(t *testing.T) {
	testCases := []struct {
		desc          string
		expected      bool
		audienceArray []string
		inputAudience string
	}{
		{
			desc:          "When audiences array contains the given audience",
			expected:      true,
			audienceArray: []string{"aud1", "aud2"},
			inputAudience: "aud1",
		},
		{
			desc:          "When audiences array not contains the given audience",
			expected:      false,
			audienceArray: []string{"aud1", "aud2"},
			inputAudience: "aud3",
		},
		{
			desc:          "When audiences array is empty",
			expected:      false,
			audienceArray: []string{},
			inputAudience: "aud3",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var config AudienceConfig
			if len(tC.audienceArray) > 0 {
				config = NewAudienceConfig(tC.audienceArray[0], tC.audienceArray[1:]...)
			} else {
				config = NewAudienceConfig("")
			}

			result := config.Contains(tC.inputAudience)

			assert.Equal(t, tC.expected, result)
		})
	}
}
