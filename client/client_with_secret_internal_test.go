package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WithSecretConfig(t *testing.T) {
	for _, testCase := range []struct {
		opts             []Option
		expectedTokenURL string
		expectedBaseURL  string
	}{
		{
			opts:             []Option{WithRealm("addons"), WithBaseURL("https://my-base.url")},
			expectedTokenURL: "https://my-base.url/addons/protocol/openid-connect/token",
			expectedBaseURL:  "https://my-base.url",
		},
		{
			opts:             []Option{WithRealm("test"), WithBaseURL("https://my-base.url")},
			expectedTokenURL: "https://my-base.url/auth/realms/test/protocol/openid-connect/token",
			expectedBaseURL:  "https://my-base.url/auth/realms",
		},
	} {

		validatorConfig := NewWithSecret("", "", WithScope(""), testCase.opts...).(*WithSecret)

		assert.Equal(t, testCase.expectedTokenURL, validatorConfig.credentials.TokenURL)
		assert.Equal(t, testCase.expectedBaseURL, validatorConfig.baseURL)
	}
}
