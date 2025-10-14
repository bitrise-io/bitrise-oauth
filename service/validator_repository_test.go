package service

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bitrise-io/bitrise-oauth/config"
	"github.com/bitrise-io/bitrise-oauth/mocks"
	"github.com/stretchr/testify/assert"
)

const (
	authServiceIssuer        = "https://auth.bitrise.io/auth/realms/bitrise-services"
	tokenIssuerServiceIssuer = "https://token-issuer.bitrise.io/auth/realms/bitrise-services"
)

func Test_GetJwtValidatorForRawToken_GivenMatchingValidatorExists_ReturnsValidator(t *testing.T) {
	authServiceValidator := NewValidator(
		config.NewAudienceConfig("bitrise-api", "bitrise"),
		WithRealm("bitrise-services"))
	tokenIssuerServiceValidator := NewValidator(
		config.NewAudienceConfig("bitrise-api", "bitrise"),
		WithRealm("bitrise-services"))
	vr := NewJwtValidatorRepository(map[string]Validator{
		authServiceIssuer:        authServiceValidator,
		tokenIssuerServiceIssuer: tokenIssuerServiceValidator,
	})

	v, err := vr.GetJwtValidatorForRawToken(mocks.RawMockToken)
	assert.NoError(t, err)

	assert.Equal(t, tokenIssuerServiceValidator, v)
}

func Test_GetJwtValidatorForRequest_GivenMatchingValidatorExists_ReturnsValidator(t *testing.T) {
	authServiceValidator := NewValidator(
		config.NewAudienceConfig("bitrise-api", "bitrise"),
		WithRealm("bitrise-services"))
	tokenIssuerServiceValidator := NewValidator(
		config.NewAudienceConfig("bitrise-api", "bitrise"),
		WithRealm("bitrise-services"))
	vr := NewJwtValidatorRepository(map[string]Validator{
		authServiceIssuer:        authServiceValidator,
		tokenIssuerServiceIssuer: tokenIssuerServiceValidator,
	})

	request, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mocks.RawMockToken))

	v, err := vr.GetJwtValidatorForRequest(request)
	assert.NoError(t, err)

	assert.Equal(t, tokenIssuerServiceValidator, v)
}

func Test_GetJwtValidatorForRequest_GivenNoMatchingValidatorExists_ReturnsError(t *testing.T) {
	authServiceValidator := NewValidator(
		config.NewAudienceConfig("bitrise-api", "bitrise"),
		WithRealm("bitrise-services"))
	vr := NewJwtValidatorRepository(map[string]Validator{
		authServiceIssuer: authServiceValidator,
	})

	request, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mocks.RawMockToken))

	_, err = vr.GetJwtValidatorForRequest(request)
	assert.EqualError(t, err, "there is no JWT validator for issuer: https://token-issuer.bitrise.io/auth/realms/bitrise-services")
}

func Test_GetJwtValidatorForRequest_GivenInvalidAuthorizationHeader_ReturnsError(t *testing.T) {
	authServiceValidator := NewValidator(
		config.NewAudienceConfig("bitrise-api", "bitrise"),
		WithRealm("bitrise-services"))
	vr := NewJwtValidatorRepository(map[string]Validator{
		authServiceIssuer: authServiceValidator,
	})

	request, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err)
	request.Header.Add("Authorization", "InvalidHeader")

	_, err = vr.GetJwtValidatorForRequest(request)
	assert.EqualError(t, err, "failed to read JWT from header")
}
