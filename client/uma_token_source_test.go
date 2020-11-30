package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
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
	responseBodyJSONString = `{"upgraded":false,"access_token":"eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwcW84eW1acnpkakdsYi1neXNyd1ByNDlDS1RlQ1dpY2tkTENsVExoZmtNIn0.eyJleHAiOjE2MDY0ODg1NDUsImlhdCI6MTYwNjQ4NDk0NSwianRpIjoiN2E3YzhlYzAtZWMyNS00MGVlLWE4NDAtNmQ0MGI3ZDE2ZjQ1IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL2F1dGgvcmVhbG1zL2JpdHJpc2Utc2VydmljZXMiLCJhdWQiOiJidWlsZC1sb2ctc2VydmljZSIsInN1YiI6IjU2M2Y2NTIyLWRlNzAtNGEyYi04YTg0LWNiYjM0NjZhNjhmNyIsInR5cCI6IkJlYXJlciIsImF6cCI6ImJ1aWxkLWxvZy1zZXJ2aWNlIiwic2Vzc2lvbl9zdGF0ZSI6IjI1MDY1M2VhLWU0MjYtNDQyNC1iZTU3LWFjNGEzM2QwM2M2ZCIsImFjciI6IjEiLCJhdXRob3JpemF0aW9uIjp7InBlcm1pc3Npb25zIjpbeyJzY29wZXMiOlsicmVhZCJdLCJjbGFpbXMiOnsidGVzdCI6WyJxd2UiLCJhc2QiXSwicXdlIjpbInF3ZTEiLCJhc2QxIl19LCJyc2lkIjoiNWE0ZWEzNDMtNGNhYi00YzhiLWE3MzMtY2ViNWZjODEzNTRhIiwicnNuYW1lIjoiYnVpbGRzIn1dfSwic2NvcGUiOiJidWlsZDpsb2ciLCJjbGllbnRIb3N0IjoiMTI3LjAuMC4xIiwiY2xpZW50SWQiOiJidWlsZC1sb2ctc2VydmljZSIsImNsaWVudEFkZHJlc3MiOiIxMjcuMC4wLjEifQ.dOlzz_HC_AppaqnAsX7JJ7Nvr-w14q34fJJhfvnptL5l_hkX4bOEiOBJeNywhkCR9o_u93I9tDDT1bIuVT2oBrn41-R_2ZRRb5xKwskPPLLnSiI0ppnpjclKkLA2fuYHw5GzMFxcvh82Uf8ZWVZlakyCGYPd9qNNk_WG-z-dFuk7sE6V--3RE2yyNV0Ei4eOhuyHa2s0EO5luLaDF3kQ2lkhdh_hBHseydmH7Za_xcySE-ogD4xfk8wUHIIlpIOpR_MDpTGIl5SZOOo0z-3r3noDp4xCGE9XEwtJSTowWuU12yzNIMoGbovUa6-DqCeoLuxrVKBwsVUVxPVGxZnBeA","expires_in":3600,"refresh_expires_in":1800,"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJlOTIyZDY0NC1hOTg1LTQ3NzAtYTI4ZS1lNjk3YTgzY2YwNzgifQ.eyJleHAiOjE2MDY0ODY3NDUsImlhdCI6MTYwNjQ4NDk0NSwianRpIjoiNGQxYzM0NTYtYjU0YS00OWUyLTlhNTktZDgxY2E4MmUxM2Q2IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL2F1dGgvcmVhbG1zL2JpdHJpc2Utc2VydmljZXMiLCJhdWQiOiJodHRwOi8vbG9jYWxob3N0OjgwODAvYXV0aC9yZWFsbXMvYml0cmlzZS1zZXJ2aWNlcyIsInN1YiI6IjU2M2Y2NTIyLWRlNzAtNGEyYi04YTg0LWNiYjM0NjZhNjhmNyIsInR5cCI6IlJlZnJlc2giLCJhenAiOiJidWlsZC1sb2ctc2VydmljZSIsInNlc3Npb25fc3RhdGUiOiIyNTA2NTNlYS1lNDI2LTQ0MjQtYmU1Ny1hYzRhMzNkMDNjNmQiLCJhdXRob3JpemF0aW9uIjp7InBlcm1pc3Npb25zIjpbeyJzY29wZXMiOlsicmVhZCJdLCJjbGFpbXMiOnsidGVzdCI6WyJxd2UiLCJhc2QiXSwicXdlIjpbInF3ZTEiLCJhc2QxIl19LCJyc2lkIjoiNWE0ZWEzNDMtNGNhYi00YzhiLWE3MzMtY2ViNWZjODEzNTRhIiwicnNuYW1lIjoiYnVpbGRzIn1dfSwic2NvcGUiOiJidWlsZDpsb2cifQ.WC1ooq49-acbnr4Ayk1UaOnQzLk3yKLjg3HEB4oNIwg","token_type":"Bearer","not-before-policy":0}`
)

var (
	testClaim = TestClaim{
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
		AccessToken:  "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwcW84eW1acnpkakdsYi1neXNyd1ByNDlDS1RlQ1dpY2tkTENsVExoZmtNIn0.eyJleHAiOjE2MDY0ODg1NDUsImlhdCI6MTYwNjQ4NDk0NSwianRpIjoiN2E3YzhlYzAtZWMyNS00MGVlLWE4NDAtNmQ0MGI3ZDE2ZjQ1IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL2F1dGgvcmVhbG1zL2JpdHJpc2Utc2VydmljZXMiLCJhdWQiOiJidWlsZC1sb2ctc2VydmljZSIsInN1YiI6IjU2M2Y2NTIyLWRlNzAtNGEyYi04YTg0LWNiYjM0NjZhNjhmNyIsInR5cCI6IkJlYXJlciIsImF6cCI6ImJ1aWxkLWxvZy1zZXJ2aWNlIiwic2Vzc2lvbl9zdGF0ZSI6IjI1MDY1M2VhLWU0MjYtNDQyNC1iZTU3LWFjNGEzM2QwM2M2ZCIsImFjciI6IjEiLCJhdXRob3JpemF0aW9uIjp7InBlcm1pc3Npb25zIjpbeyJzY29wZXMiOlsicmVhZCJdLCJjbGFpbXMiOnsidGVzdCI6WyJxd2UiLCJhc2QiXSwicXdlIjpbInF3ZTEiLCJhc2QxIl19LCJyc2lkIjoiNWE0ZWEzNDMtNGNhYi00YzhiLWE3MzMtY2ViNWZjODEzNTRhIiwicnNuYW1lIjoiYnVpbGRzIn1dfSwic2NvcGUiOiJidWlsZDpsb2ciLCJjbGllbnRIb3N0IjoiMTI3LjAuMC4xIiwiY2xpZW50SWQiOiJidWlsZC1sb2ctc2VydmljZSIsImNsaWVudEFkZHJlc3MiOiIxMjcuMC4wLjEifQ.dOlzz_HC_AppaqnAsX7JJ7Nvr-w14q34fJJhfvnptL5l_hkX4bOEiOBJeNywhkCR9o_u93I9tDDT1bIuVT2oBrn41-R_2ZRRb5xKwskPPLLnSiI0ppnpjclKkLA2fuYHw5GzMFxcvh82Uf8ZWVZlakyCGYPd9qNNk_WG-z-dFuk7sE6V--3RE2yyNV0Ei4eOhuyHa2s0EO5luLaDF3kQ2lkhdh_hBHseydmH7Za_xcySE-ogD4xfk8wUHIIlpIOpR_MDpTGIl5SZOOo0z-3r3noDp4xCGE9XEwtJSTowWuU12yzNIMoGbovUa6-DqCeoLuxrVKBwsVUVxPVGxZnBeA",
		TokenType:    "Bearer",
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJlOTIyZDY0NC1hOTg1LTQ3NzAtYTI4ZS1lNjk3YTgzY2YwNzgifQ.eyJleHAiOjE2MDY0ODY3NDUsImlhdCI6MTYwNjQ4NDk0NSwianRpIjoiNGQxYzM0NTYtYjU0YS00OWUyLTlhNTktZDgxY2E4MmUxM2Q2IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL2F1dGgvcmVhbG1zL2JpdHJpc2Utc2VydmljZXMiLCJhdWQiOiJodHRwOi8vbG9jYWxob3N0OjgwODAvYXV0aC9yZWFsbXMvYml0cmlzZS1zZXJ2aWNlcyIsInN1YiI6IjU2M2Y2NTIyLWRlNzAtNGEyYi04YTg0LWNiYjM0NjZhNjhmNyIsInR5cCI6IlJlZnJlc2giLCJhenAiOiJidWlsZC1sb2ctc2VydmljZSIsInNlc3Npb25fc3RhdGUiOiIyNTA2NTNlYS1lNDI2LTQ0MjQtYmU1Ny1hYzRhMzNkMDNjNmQiLCJhdXRob3JpemF0aW9uIjp7InBlcm1pc3Npb25zIjpbeyJzY29wZXMiOlsicmVhZCJdLCJjbGFpbXMiOnsidGVzdCI6WyJxd2UiLCJhc2QiXSwicXdlIjpbInF3ZTEiLCJhc2QxIl19LCJyc2lkIjoiNWE0ZWEzNDMtNGNhYi00YzhiLWE3MzMtY2ViNWZjODEzNTRhIiwicnNuYW1lIjoiYnVpbGRzIn1dfSwic2NvcGUiOiJidWlsZDpsb2cifQ.WC1ooq49-acbnr4Ayk1UaOnQzLk3yKLjg3HEB4oNIwg",
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

type TestClaim struct {
	Param1 []string `json:"param1"`
	Param2 []string `json:"param2"`
}

func Test_GivenClaim_WhenEncodeClaimCalled_ThenExpectTheEncodedClaimToBeReturned(t *testing.T) {
	// When
	encodedClaim, err := encodeClaim(&testClaim)

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
