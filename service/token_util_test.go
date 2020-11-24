package service

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	auth0 "github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	authorization        = "Authorization"
	bearer               = "Bearer"
	defaultRequestURL    = "http://localhost"
	defaultRequestMethod = http.MethodGet

	defaultIssuer = "issuer"
	defaultKid    = ""
)

var (
	defaultAudience       = []string{"audience"}
	defaultSecret         = genRSASSAJWK(jose.RS256, defaultKid)
	defaultSecretProvider = auth0.NewKeyProvider(defaultSecret.Public())
)

func getTestTokenWithKid(audience []string, issuer string, expTime time.Time, alg jose.SignatureAlgorithm, key interface{}, kid string) string {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: alg, Key: key}, (&jose.SignerOptions{ExtraHeaders: map[jose.HeaderKey]interface{}{"kid": kid}}).WithType("JWT"))
	if err != nil {
		panic(err)
	}

	cl := jwt.Claims{
		Issuer:   issuer,
		Audience: audience,
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Expiry:   jwt.NewNumericDate(expTime),
	}

	tokenStr, err := jwt.Signed(signer).Claims(cl).CompactSerialize()
	if err != nil {
		panic(err)
	}

	return tokenStr
}

func genRSASSAJWK(sigAlg jose.SignatureAlgorithm, kid string) jose.JSONWebKey {
	var bits int
	if sigAlg == jose.RS256 {
		bits = 2048
	}
	if sigAlg == jose.RS512 {
		bits = 4096
	}

	key, _ := rsa.GenerateKey(rand.Reader, bits)

	jsonWebKey := jose.JSONWebKey{
		Key:       key,
		KeyID:     kid,
		Use:       "sig",
		Algorithm: string(sigAlg),
	}

	return jsonWebKey
}

func createRequestWithToken(token string) *http.Request {
	request, err := http.NewRequest(defaultRequestMethod, defaultRequestURL, nil)
	if err != nil {
		panic("Can't create request.")
	}

	authHeader := fmt.Sprintf("%s %s", bearer, token)
	request.Header.Add(authorization, authHeader)
	return request
}
