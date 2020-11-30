package service

import (
	"errors"
	"testing"
	"time"

	"github.com/c2fo/testify/assert"
	"github.com/c2fo/testify/require"
	"gopkg.in/square/go-jose.v2/jwt"
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

// Helpers

func givenTokenWithClaims(claims interface{}) TokenWithClaims {
	token, _ := newTestTokenConfig().newTokenWithClaims(claims)
	return TokenWithClaims{
		key:   defaultSecret.Public(),
		token: token,
	}
}

func givenClaimsWithoutAuthorization() interface{} {
	return struct {
		Issuer   string
		Audience jwt.Audience
		IssuedAt jwt.NumericDate
		Expiry   jwt.NumericDate
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
