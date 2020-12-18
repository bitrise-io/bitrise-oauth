package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/bitrise-io/bitrise-oauth/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	expectedEncodedClaim   = "eyJwYXJhbTEiOlsidmFsdWUxIl0sInBhcmFtMiI6WyJ2YWx1ZTIiXX0="
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
	realm  = "testRealm"
)

var (
	testClaims = testClaim{
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

type testClaim struct {
	Param1 []string `json:"param1"`
	Param2 []string `json:"param2"`
}

func Test_GivenClaim_WhenEncodeClaimCalled_ThenExpectTheEncodedClaimToBeReturned(t *testing.T) {
	// When
	encodedClaim, err := encodeClaim(&testClaims)

	// Then
	require.NoError(t, err)
	assert.Equal(t, expectedEncodedClaim, encodedClaim)
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
	request, err := umaTokenSource.newTokenRequest(expectedEncodedClaim, testPermission, audienceConfig)
	b, err := ioutil.ReadAll(request.Body)
	body := string(b)

	// Then
	require.NoError(t, err)

	assert.Equal(t, formURLEncoded, request.Header[contentType][0])

	assert.Contains(t, body, urlEncodedBodyParam(grantType, umaGrantType))
	assert.Contains(t, body, urlEncodedBodyParam(claimToken, expectedEncodedClaim))
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
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseBodyJSONString)),
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
		Body:       ioutil.NopCloser(bytes.NewBufferString("responseBody")),
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

func Test_GivenUMATokenSource_WhenTokenAsked_NoErrors(t *testing.T){
	ts := newAssertingMockServer(t, realm, func(t *testing.T, r *http.Request) {

	})
	defer ts.Close()
	authProvider := NewWithSecret(
		"test-client-id",
		"test-client-secret",
		WithScope("test"),
		WithBaseURL(ts.URL),
		WithRealm(realm))
	tokenSource := authProvider.UMATokenSource()
	_, err := tokenSource.Token(nil, nil, audienceConfig)
	require.NoError(t, err)
}

func Test_GivenAudienceConfiguration_WhenUMATokenSourceIsInstantiated_ThenTokenCallsWithAudience(t *testing.T){
	cases := map[string] struct{
		AudienceFromSourceOptions string
		AudienceFromTokenOptions string
		ExpectedOptionsSent []string
	}{
		"No audience" : {
			AudienceFromTokenOptions: "",
			AudienceFromSourceOptions: "",
			ExpectedOptionsSent: []string{},
		},
		"Different audiences are provided from each" : {
			AudienceFromTokenOptions: "aud-cof-from-token-method",
			AudienceFromSourceOptions: "aud-conf-from-options",
			ExpectedOptionsSent: []string{"aud-cof-from-token-method", "aud-conf-from-options"},
		},
		"Only token option provided" : {
		AudienceFromTokenOptions: "aud-cof-from-token-method",
		AudienceFromSourceOptions: "",
		ExpectedOptionsSent: []string{"aud-cof-from-token-method"},
		},
		"Only source option provided" : {
			AudienceFromTokenOptions: "",
			AudienceFromSourceOptions: "aud-conf-from-options",
			ExpectedOptionsSent: []string{"aud-conf-from-options"},
		},
		"Both provides same audience" : {
			AudienceFromTokenOptions: "test-aud",
			AudienceFromSourceOptions: "test-aud",
			ExpectedOptionsSent: []string{"test-aud"},
		},
	}

	for _, c := range cases {
		runAudiencesTestCase(
			t,
			c.AudienceFromTokenOptions,
			c.AudienceFromSourceOptions,
			c.ExpectedOptionsSent)
	}
}

func runAudiencesTestCase(
	t *testing.T,
	audFromToken string,
	audFromOptions string,
	expectedOptionsSet []string ) {
	ts := newAssertingMockServer(t, realm, func(t *testing.T, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		body := string(b)
		for _, expected := range expectedOptionsSet {
			assert.Contains(t, body, urlEncodedBodyParam(audience, expected))
		}
	})
	defer ts.Close()

	authProvider := NewWithSecret(
		"test-client-id",
		"test-client-secret",
		WithScope("test"),
		WithBaseURL(ts.URL),
		WithRealm(realm))
	audConfOpt := WithAudienceConfig(config.NewAudienceConfig(audFromOptions))
	tokenSource := authProvider.UMATokenSource(audConfOpt)
	_ ,err := tokenSource.Token(nil, nil, config.NewAudienceConfig(audFromToken))
	require.NoError(t, err)
}

func newAssertingMockServer(
	t *testing.T,
	realm string,
	assertFunc func(t *testing.T,  r *http.Request)) *httptest.Server{
	tokenEndpointURL := "/auth/realms/" + realm +"/protocol/openid-connect/token"
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case tokenEndpointURL:
			assertFunc(t, r)
			json.NewEncoder(w).Encode(tokenJSON{
				AccessToken:  "my-test-token",
				RefreshToken: "refresh-token",
				TokenType:    "Bearer",
				ExpiresIn:    3600, 
			})
			w.WriteHeader(http.StatusOK)
		}
	}))
}