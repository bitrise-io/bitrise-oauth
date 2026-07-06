package oauthtest_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bitrise-io/bitrise-oauth/config"
	"github.com/bitrise-io/bitrise-oauth/service"
	"github.com/bitrise-io/bitrise-oauth/service/oauthtest"
)

const (
	testRealm    = "bitrise-services"
	testAudience = "build-log-service"
	testResource = "builds"
)

type buildClaims struct {
	AppSlug []string `json:"app_slug"`
	BuildID []string `json:"build_id"`
}

func newValidator(srv *oauthtest.Server) service.Validator {
	return service.NewValidator(
		config.NewAudienceConfig(testAudience),
		service.WithBaseURL(srv.BaseURL()),
		service.WithRealm(testRealm),
		service.WithJWKSUrl(srv.JWKSURL()),
	)
}

func requestWithToken(token string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost/anything", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// TestMintedTokenPassesValidator is the anti-drift guarantee: a token minted by
// oauthtest must validate through the real service.Validator and expose its
// permissions and claims. If the validator's expected shape changes, this fails.
func TestMintedTokenPassesValidator(t *testing.T) {
	srv, err := oauthtest.Start(testRealm)
	require.NoError(t, err)
	defer srv.Close()

	token, err := srv.Token(oauthtest.TokenOptions{
		Audience: []string{testAudience},
		Permissions: []oauthtest.Permission{{
			Resource: testResource,
			Scopes:   []string{"write"},
			Claims:   buildClaims{AppSlug: []string{"app-1"}, BuildID: []string{"build-1"}},
		}},
	})
	require.NoError(t, err)

	twc, err := newValidator(srv).ValidateRequestAndReturnToken(requestWithToken(token))
	require.NoError(t, err)

	require.NoError(t, twc.ValidatePermissionScopes(testResource, []string{"write"}))

	var got buildClaims
	require.NoError(t, twc.Claim(testResource, &got))
	assert.Equal(t, []string{"app-1"}, got.AppSlug)
	assert.Equal(t, []string{"build-1"}, got.BuildID)
}

func TestWrongAudienceRejected(t *testing.T) {
	srv, err := oauthtest.Start(testRealm)
	require.NoError(t, err)
	defer srv.Close()

	token, err := srv.Token(oauthtest.TokenOptions{
		Audience: []string{"some-other-service"},
		Permissions: []oauthtest.Permission{{
			Resource: testResource,
			Scopes:   []string{"write"},
		}},
	})
	require.NoError(t, err)

	err = newValidator(srv).ValidateRequest(requestWithToken(token))
	assert.Error(t, err)
}

func TestTokenSignedByDifferentKeyRejected(t *testing.T) {
	srv, err := oauthtest.Start(testRealm)
	require.NoError(t, err)
	defer srv.Close()

	// A second server has a different key pair; its token must not validate
	// against the first server's JWKS.
	other, err := oauthtest.Start(testRealm)
	require.NoError(t, err)
	defer other.Close()

	token, err := other.Token(oauthtest.TokenOptions{
		Audience: []string{testAudience},
		Issuer:   srv.Issuer(),
	})
	require.NoError(t, err)

	err = newValidator(srv).ValidateRequest(requestWithToken(token))
	assert.Error(t, err)
}
