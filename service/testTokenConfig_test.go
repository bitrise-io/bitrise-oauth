package service

// import (
// 	"crypto/rand"
// 	"crypto/rsa"
// 	"fmt"
// 	"net/http"
// 	"testing"
// 	"time"

// 	auth0 "github.com/auth0-community/go-auth0"
// 	"github.com/stretchr/testify/require"
// 	"gopkg.in/square/go-jose.v2"
// 	"gopkg.in/square/go-jose.v2/jwt"
// )

// const (
// 	authorization        = "Authorization"
// 	bearer               = "Bearer"
// 	defaultRequestURL    = "http://localhost"
// 	defaultRequestMethod = http.MethodGet

// 	defaultIssuer = "issuer"
// 	defaultKid    = ""
// )

// var (
// 	defaultAudience       = []string{"audience"}
// 	defaultSecret         = genRSASSAJWK(jose.RS256, defaultKid)
// 	defaultSecretProvider = auth0.NewKeyProvider(defaultSecret.Public())
// )

// type testTokenConfig struct {
// 	audience []string
// 	issuer   string
// 	expTime  time.Time
// 	alg      jose.SignatureAlgorithm
// 	key      interface{}
// 	kid      string
// }

// func newTestToken() testTokenConfig {
// 	return testTokenConfig{
// 		audience: defaultAudience,
// 		issuer:   defaultIssuer,
// 		expTime:  time.Now().Add(24 * time.Hour),
// 		alg:      jose.RS256,
// 		key:      defaultSecret,
// 		kid:      defaultKid,
// 	}
// }

// func (testTokenConfig testTokenConfig) getTokenStringWithKid(t *testing.T) string {
// 	tokenStr, err := testTokenConfig.getJWTBuilder(t).CompactSerialize()
// 	if err != nil {
// 		require.Fail(t, err.Error())
// 	}

// 	return tokenStr
// }

// func (testTokenConfig testTokenConfig) getTokenWithKid(t *testing.T) *jwt.JSONWebToken {
// 	token, err := testTokenConfig.getJWTBuilder(t).Token()
// 	if err != nil {
// 		require.Fail(t, err.Error())
// 	}

// 	return token
// }

// func (testTokenConfig testTokenConfig) getJWTBuilder(t *testing.T) jwt.Builder {
// 	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: testTokenConfig.alg, Key: testTokenConfig.key}, (&jose.SignerOptions{ExtraHeaders: map[jose.HeaderKey]interface{}{"kid": testTokenConfig.kid}}).WithType("JWT"))
// 	if err != nil {
// 		require.Fail(t, err.Error())
// 	}

// 	cl := jwt.Claims{
// 		Issuer:   testTokenConfig.issuer,
// 		Audience: testTokenConfig.audience,
// 		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
// 		Expiry:   jwt.NewNumericDate(testTokenConfig.expTime),
// 	}

// 	return jwt.Signed(signer).Claims(cl)
// }

// func (testTokenConfig testTokenConfig) createRequest(t *testing.T) *http.Request {
// 	request, err := http.NewRequest(defaultRequestMethod, defaultRequestURL, nil)
// 	if err != nil {
// 		require.Fail(t, "Can't create request.")
// 	}

// 	authHeader := fmt.Sprintf("%s %s", bearer, testTokenConfig.getTokenStringWithKid(t))
// 	request.Header.Add(authorization, authHeader)
// 	return request
// }

// func genRSASSAJWK(sigAlg jose.SignatureAlgorithm, kid string) jose.JSONWebKey {
// 	var bits int
// 	if sigAlg == jose.RS256 {
// 		bits = 2048
// 	}
// 	if sigAlg == jose.RS512 {
// 		bits = 4096
// 	}

// 	key, err := rsa.GenerateKey(rand.Reader, bits)
// 	if err != nil {
// 		panic(err)
// 	}

// 	jsonWebKey := jose.JSONWebKey{
// 		Key:       key,
// 		KeyID:     kid,
// 		Use:       "sig",
// 		Algorithm: string(sigAlg),
// 	}

// 	return jsonWebKey
// }
