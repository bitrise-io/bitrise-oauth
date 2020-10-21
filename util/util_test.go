package util_test

import (
	"testing"

	"github.com/bitrise-io/bitrise-oauth/util"
	"github.com/c2fo/testify/assert"
)

const (
	expectedURL = "https://auth.services.bitrise.io/auth/realms/master"
)

func Test_GivenBaseUrlAndPathPieces_WhenJoinUrlCalled_ThenExpectTheJoinedUrl(t *testing.T) {
	testCases := []struct {
		name       string
		baseURL    string
		pathPieces []string
		want       string
	}{
		{
			"1. Base url without traling /, first path piece without leading /",
			"https://auth.services.bitrise.io",
			[]string{"auth/realms", "master"},
			expectedURL,
		},
		{
			"2. Base url with traling /, first path piece without leading /",
			"https://auth.services.bitrise.io/",
			[]string{"auth/realms", "master"},
			expectedURL,
		},
		{
			"3. First path piece without leading /, first path piece with leading /",
			"https://auth.services.bitrise.io",
			[]string{"/auth/realms", "master"},
			expectedURL,
		},
		{
			"4. First path piece with leading /, first path piece with leading /",
			"https://auth.services.bitrise.io/",
			[]string{"/auth/realms", "master"},
			expectedURL,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// When
			actual := util.JoinURL(testCase.baseURL, testCase.pathPieces...)

			// Then
			assert.Equal(t, testCase.want, actual)
		})
	}
}
