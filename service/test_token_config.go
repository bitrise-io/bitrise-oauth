package service

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/bitrise-io/go-auth0"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const (
	authorizationHeader  = "Authorization"
	bearer               = "Bearer"
	defaultRequestURL    = "http://localhost"
	defaultRequestMethod = http.MethodGet

	defaultIssuer = "issuer"
	defaultKid    = "kid"
)

var (
	defaultAudience       = []string{"audience"}
	defaultSecret         = genRSASSAJWK(jose.RS256, defaultKid)
	defaultSecretProvider = auth0.NewKeyProvider(defaultSecret.Public())
)

type testTokenConfig struct {
	audience []string
	issuer   string
	expTime  time.Time
	alg      jose.SignatureAlgorithm
	key      interface{}
	kid      string
}

func newTestTokenConfig() testTokenConfig {
	return testTokenConfig{
		defaultAudience,
		defaultIssuer,
		time.Now().Add(24 * time.Hour),
		jose.RS256,
		defaultSecret,
		defaultKid,
	}
}

func (testToken testTokenConfig) getTokenString() string {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: testToken.alg, Key: testToken.key}, (&jose.SignerOptions{ExtraHeaders: map[jose.HeaderKey]interface{}{"kid": testToken.kid}}).WithType("JWT"))
	if err != nil {
		panic(err)
	}

	cl := jwt.Claims{
		Issuer:   testToken.issuer,
		Audience: testToken.audience,
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Expiry:   jwt.NewNumericDate(testToken.expTime),
	}

	tokenStr, err := jwt.Signed(signer).Claims(cl).Serialize()
	if err != nil {
		panic(err)
	}

	return tokenStr
}

func (testToken testTokenConfig) newTokenWithClaims(claims interface{}) (*jwt.JSONWebToken, interface{}) {
	actKey := jose.SigningKey{Algorithm: testToken.alg, Key: testToken.key}
	signer, err := jose.NewSigner(actKey, (&jose.SignerOptions{ExtraHeaders: map[jose.HeaderKey]interface{}{"kid": testToken.kid}}).WithType("JWT"))
	if err != nil {
		panic(err)
	}

	token, err := jwt.Signed(signer).Claims(claims).Token()
	if err != nil {
		panic(err)
	}

	return token, actKey
}

func genRSASSAJWK(sigAlg jose.SignatureAlgorithm, kid string) jose.JSONWebKey {
	var bits int
	if sigAlg == jose.RS256 {
		bits = 2048
	}
	if sigAlg == jose.RS512 {
		bits = 4096
	}

	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}

	jsonWebKey := jose.JSONWebKey{
		Key:       key,
		KeyID:     kid,
		Use:       "sig",
		Algorithm: string(sigAlg),
	}

	return jsonWebKey
}

func (testToken testTokenConfig) newRequest() *http.Request {
	request, err := http.NewRequest(defaultRequestMethod, defaultRequestURL, nil)
	if err != nil {
		panic("Can't create request.")
	}

	tokenString := testToken.getTokenString()
	authHeader := fmt.Sprintf("%s %s", bearer, tokenString)
	request.Header.Add(authorizationHeader, authHeader)
	return request
}
