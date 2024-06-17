package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bitrise-io/bitrise-oauth/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
)

const (
	expectedEncodedPayload = "eyJwYXJhbTEiOlsidmFsdWUxIl0sInBhcmFtMiI6WyJ2YWx1ZTIiXX0="
	defaultAudience        = "audience"
	responseBodyJSONString = `
	{
		"upgraded":false,
		"access_token":"access_token",
		"expires_in":3600,
		"refresh_expires_in":1800,
		"refresh_token":"refresh_token",
		"token_type":"Bearer",
		"not-before-policy":0
	}`
)

var (
	testPayloads = testPayload{
		Param1: []string{"value1"},
		Param2: []string{"value2"},
	}

	testPermission = []Permission{Permission{
		resourceName:       "resourceName",
		authorizationScope: "scope",
	}}

	audienceConfig = config.NewAudienceConfig(defaultAudience)

	// investigate shorter tokens
	tj = tokenJSON{
		AccessToken:  "access_token",
		TokenType:    "Bearer",
		RefreshToken: "refresh_token",
		ExpiresIn:    3600,
	}

	baseTime = time.Now()

	testToken = &oauth2.Token{
		AccessToken:  tj.AccessToken,
		TokenType:    tj.TokenType,
		RefreshToken: tj.RefreshToken,
		Expiry:       tj.expiry(baseTime),
	}
)

type testPayload struct {
	Param1 []string `json:"param1"`
	Param2 []string `json:"param2"`
}

func Test_GivenPayload_WhenEncodePayloadCalled_ThenExpectTheEncodedPayloadToBeReturned(t *testing.T) {
	// When
	encodedPayload, err := encodePayload(&testPayloads)

	// Then
	require.NoError(t, err)
	assert.Equal(t, expectedEncodedPayload, encodedPayload)
}

func Test_GivenTokenSource_WhenATokenRequestIsCreated_ThenExpectParamsToBeOnTheRequest(t *testing.T) {
	// Given
	umaTokenSource := newUMATokenSource(clientcredentials.Config{
		ClientID:     "clientID",
		ClientSecret: "clientSecret",
		TokenURL:     "tokenURL",
		Scopes:       []string{"scope"},
	})

	// When
	request, err := umaTokenSource.newTokenRequest(expectedEncodedPayload, testPermission, audienceConfig)
	require.NoError(t, err)
	b, err := io.ReadAll(request.Body)
	body := string(b)

	// Then
	require.NoError(t, err)

	assert.Equal(t, formURLEncoded, request.Header[contentType][0])

	assert.Contains(t, body, urlEncodedBodyParam(grantType, umaGrantType))
	assert.Contains(t, body, urlEncodedBodyParam(claimToken, expectedEncodedPayload))
	assert.Contains(t, body, urlEncodedBodyParam(claimTokenFormat, umaClaimTokenFormat))
	assert.Contains(t, body, urlEncodedBodyParam(clientID, umaTokenSource.config.ClientID))
	assert.Contains(t, body, urlEncodedBodyParam(clientSecret, umaTokenSource.config.ClientSecret))
	assert.Contains(t, body, urlEncodedBodyParam(permission, testPermission[0].requestParam()))
	assert.Contains(t, body, urlEncodedBodyParam(audience, audienceConfig.All()[0]))
}

func Test_GivenSuccessfulTokenResponse_WhenBodyIsExtracted_ThenExpectTheBodyBeReturned(t *testing.T) {
	// Given
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(responseBodyJSONString)),
	}

	// When
	body, err := extractResponseBody(response)

	// Then
	require.NoError(t, err)
	assert.Equal(t, responseBodyJSONString, string(body))
}

func Test_GivenUnsuccessfulTokenResponse_WhenBodyIsExtracted_ThenExpectAnError(t *testing.T) {
	// Given
	response := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(bytes.NewBufferString("responseBody")),
	}

	// When
	_, err := extractResponseBody(response)

	// Then
	require.Error(t, err)
}

func Test_GivenSuccessfulTokenResponse_WhenTokenIsExtractedFromBody_ThenExpectTokenToBeReturned(t *testing.T) {
	// When
	token, err := extractTokenFromBody([]byte(responseBodyJSONString), baseTime)

	// Then
	require.NoError(t, err)
	assert.Equal(t, token, testToken)
}

func urlEncodedBodyParam(key, value string) string {
	return fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value))
}
