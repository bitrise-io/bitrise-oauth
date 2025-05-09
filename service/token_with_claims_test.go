package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/go-jose/go-jose.v2/jwt"
)

// Payload
func Test_Payload(t *testing.T) {
	// Given
	tokenWithClaims := givenTokenWithClaims(givenClaims())

	// When
	data, err := tokenWithClaims.Payload()

	// Then
	require.NoError(t, err)
	assert.NotEmpty(t, data["iss"])
	assert.NotEmpty(t, data["aud"])
	assert.NotEmpty(t, data["iat"])
	assert.NotEmpty(t, data["exp"])
	assert.NotEmpty(t, data["authorization"])
}

// Permissions
func Test_Permissions(t *testing.T) {
	// Given
	tokenWithClaims := givenTokenWithClaims(givenClaims())

	// When
	permissions, err := tokenWithClaims.Permissions()

	// Then
	require.NoError(t, err)
	assert.NotNil(t, permissions)
}

func Test_PermissionsError(t *testing.T) {
	testCases := []struct {
		name   string
		claims interface{}
		want   error
	}{
		{
			"1. Given token missing persmissons then expect error",
			givenClaimsWithoutPermissons(),
			errors.New("permissions is missing from token"),
		},
		{
			"2. Given token missing authorization then expect error",
			givenClaimsWithoutAuthorization(),
			errors.New("authorization is missing from token"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Given
			tokenWithClaims := givenTokenWithClaims(testCase.claims)

			// When
			_, err := tokenWithClaims.Permissions()

			// Then
			assert.EqualError(t, err, testCase.want.Error())
		})
	}
}

// Claims
type testBuildClaim struct {
	BuildID string `json:"build_id"`
}

type testScopeClaim struct {
	Scope string `json:"scope,omitempty"`
}

func Test_GivenResourceWithClaim_WhenClaimsCalledWithValidResourceName_ThenExpectClaimReturned(t *testing.T) {
	// Given
	resourceName := "resourceName"
	expectedBuildClaim := testBuildClaim{
		BuildID: "test_build_id",
	}
	claims := givenClaimsWithResource(resourceName, expectedBuildClaim)
	tokenWithClaims := givenTokenWithClaims(claims)

	// When
	actualBuildClaim := testBuildClaim{}
	err := tokenWithClaims.Claim(resourceName, &actualBuildClaim)

	// Then
	require.NoError(t, err)
	assert.Equal(t, expectedBuildClaim, actualBuildClaim)
}

func Test_GivenResourceWithClaim_WhenClaimsCalledWithInvalidResourceName_ThenExpectError(t *testing.T) {
	// Given
	buildClaim := testBuildClaim{
		BuildID: "test_build_id",
	}
	claims := givenClaimsWithResource("whatever", buildClaim)
	tokenWithClaims := givenTokenWithClaims(claims)

	// When
	parsedBuildClaim := testBuildClaim{}
	err := tokenWithClaims.Claim("invalid_resource_name", &parsedBuildClaim)

	// Then
	assert.EqualError(t, err, "permission for resource: invalid_resource_name not found")
}

func Test_ValidateScopes_WhenScopeClaimIsMissing_ThenExpectError(t *testing.T) {
	// Given
	tokenWithClaims := givenTokenWithClaims(struct{}{})

	// When
	err := tokenWithClaims.ValidateScopes([]string{"app:read", "missing:write"})

	// Then
	assert.EqualError(t, err, "no scope claim in token")
}

func Test_ValidateScopes_WhenGivenScopeIsMissing_ThenExpectError(t *testing.T) {
	// Given
	scopeClaim := testScopeClaim{
		Scope: "app:read build:write",
	}
	tokenWithClaims := givenTokenWithClaims(scopeClaim)

	// When
	err := tokenWithClaims.ValidateScopes([]string{"app:read", "missing:write"})

	// Then
	assert.EqualError(t, err, "scope missing:write is missing from the token")
}

func Test_ValidateScopes_WhenAllScopesAreFound_ThenExpectNoError(t *testing.T) {
	// Given
	scopeClaim := testScopeClaim{
		Scope: "app:read build:write",
	}
	tokenWithClaims := givenTokenWithClaims(scopeClaim)

	// When
	err1 := tokenWithClaims.ValidateScopes([]string{"build:write", "app:read"})
	err2 := tokenWithClaims.ValidateScopes([]string{"build:write", "app:read"})

	// Then
	require.NoError(t, err1)
	require.NoError(t, err2)
}

func Test_ValidatePermissionScopes_WhenResourceNameIsInvalid_ThenExpectError(t *testing.T) {
	// Given
	claims := givenClaimsWithAuthorization(
		authorization{
			Permissions: []permisson{
				{
					Scopes: nil,
					Claims: nil,
					Rsid:   "test-id",
					Rsname: "builds",
				},
			},
		},
	)
	tokenWithClaims := givenTokenWithClaims(claims)

	// When
	err := tokenWithClaims.ValidatePermissionScopes("invalid-resource-name", []string{"write"})

	// Then
	assert.EqualError(t, err, "resource name invalid-resource-name does not match with any resources in the token")
}

func Test_ValidatePermissionScopes_WhenScopesAreMissing_ThenExpectError(t *testing.T) {
	// Given
	claims := givenClaimsWithAuthorization(
		authorization{
			Permissions: []permisson{
				{
					Scopes: nil,
					Claims: nil,
					Rsid:   "test-id",
					Rsname: "builds",
				},
			},
		},
	)
	tokenWithClaims := givenTokenWithClaims(claims)

	// When
	err := tokenWithClaims.ValidatePermissionScopes("builds", []string{"write"})

	// Then
	assert.EqualError(t, err, "no permission scope claim in token")
}

func Test_ValidatePermissionScopes_WhenGivenScopeIsMissing_ThenExpectError(t *testing.T) {
	// Given
	claims := givenClaimsWithAuthorization(
		authorization{
			Permissions: []permisson{
				{
					Scopes: []string{"read"},
					Claims: nil,
					Rsid:   "test-id",
					Rsname: "builds",
				},
			},
		},
	)
	tokenWithClaims := givenTokenWithClaims(claims)

	// When
	err := tokenWithClaims.ValidatePermissionScopes("builds", []string{"write"})

	// Then
	assert.EqualError(t, err, "scope write is missing from permissions")
}

func Test_ValidatePermissionScopes_WhenAllScopesAreFound_ThenExpectError(t *testing.T) {
	// Given
	claims := givenClaimsWithAuthorization(
		authorization{
			Permissions: []permisson{
				{
					Scopes: []string{"read", "write"},
					Claims: nil,
					Rsid:   "test-id",
					Rsname: "builds",
				},
			},
		},
	)
	tokenWithClaims := givenTokenWithClaims(claims)

	// When
	err := tokenWithClaims.ValidatePermissionScopes("builds", []string{"read", "write"})

	// Then
	require.NoError(t, err)
}

// Helpers

func givenTokenWithClaims(claims interface{}) tokenWithClaims {
	token, _ := newTestTokenConfig().newTokenWithClaims(claims)
	return tokenWithClaims{
		key:   defaultSecret.Public(),
		token: token,
	}
}

func givenClaimsWithoutAuthorization() interface{} {
	return struct {
		Issuer   string
		Audience jwt.Audience
		IssuedAt *jwt.NumericDate
		Expiry   *jwt.NumericDate
	}{
		Issuer:   defaultIssuer,
		Audience: defaultAudience,
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Expiry:   jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
}

func givenClaims() umaToken {
	return givenClaimsWithResource("whatever", nil)
}

func givenClaimsWithoutPermissons() umaToken {
	return givenClaimsWithAuthorization(authorization{})
}

func givenClaimsWithResource(resourceName string, claim interface{}) umaToken {
	auth := authorization{
		Permissions: []permisson{
			{
				Scopes: []string{"scope"},
				Claims: claim,
				Rsid:   "resourceId",
				Rsname: resourceName,
			},
		},
	}

	return givenClaimsWithAuthorization(auth)
}

func givenClaimsWithAuthorization(auth authorization) umaToken {
	return umaToken{
		Issuer:        defaultIssuer,
		Audience:      defaultAudience,
		IssuedAt:      jwt.NewNumericDate(time.Now().UTC()),
		Expiry:        jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Authorization: auth,
	}
}
