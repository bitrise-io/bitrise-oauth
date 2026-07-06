// Package oauthtest provides an in-process JWKS server and UMA token minter for
// exercising the bitrise-oauth service validator end-to-end in tests, without
// depending on a live auth service or Token Issuer Service.
//
// A Server generates an RSA key pair, serves the public half as a JWKS, and
// mints tokens signed with the private half. Point a service.Validator (or a
// whole service process via its TOKEN_ISSUER_SERVICE_BASE_URL) at the Server
// and the tokens it issues validate against the real validation path: JWKS
// fetch, signature verification, issuer and audience checks.
package oauthtest

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const defaultKeyID = "oauthtest-key"

// Server serves a JWKS for a generated RSA key pair and mints UMA tokens signed
// with the matching private key.
type Server struct {
	realm   string
	kid     string
	key     *rsa.PrivateKey
	jwk     jose.JSONWebKey
	srv     *http.Server
	baseURL string
}

// Option configures a Server.
type Option func(*Server)

// WithKeyID overrides the key ID advertised in the JWKS and the token header.
func WithKeyID(kid string) Option {
	return func(s *Server) { s.kid = kid }
}

// Start generates an RSA key pair and starts an HTTP server serving the JWKS.
// It serves the key set at both the Keycloak-style path (the validator default,
// …/protocol/openid-connect/certs) and the shorter TIS-style path (…/certs), so
// it works regardless of how the validator's JWKS URL is configured. Call Close
// when done.
func Start(realm string, opts ...Option) (*Server, error) {
	s := &Server{realm: realm, kid: defaultKeyID}
	for _, opt := range opts {
		opt(s)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generate RSA key: %w", err)
	}
	s.key = key
	s.jwk = jose.JSONWebKey{
		Key:       key,
		KeyID:     s.kid,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}
	s.baseURL = "http://" + ln.Addr().String()

	realmPath := fmt.Sprintf("/auth/realms/%s", realm)
	mux := http.NewServeMux()
	mux.HandleFunc(realmPath+"/certs", s.handleJWKS)
	mux.HandleFunc(realmPath+"/protocol/openid-connect/certs", s.handleJWKS)

	s.srv = &http.Server{Handler: mux}
	go func() { _ = s.srv.Serve(ln) }()

	return s, nil
}

// Close stops the server.
func (s *Server) Close() error { return s.srv.Close() }

// BaseURL is the server root. Feed it to the service under test as its auth
// base URL (e.g. TOKEN_ISSUER_SERVICE_BASE_URL).
func (s *Server) BaseURL() string { return s.baseURL }

// Issuer is the token issuer, matching what the validator derives from the realm.
func (s *Server) Issuer() string {
	return fmt.Sprintf("%s/auth/realms/%s", s.baseURL, s.realm)
}

// JWKSURL is the TIS-style JWKS endpoint. Pass to service.WithJWKSUrl when the
// validator is configured with an explicit JWKS URL.
func (s *Server) JWKSURL() string {
	return fmt.Sprintf("%s/auth/realms/%s/certs", s.baseURL, s.realm)
}

func (s *Server) handleJWKS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{s.jwk.Public()}})
}

// Permission is one UMA authorization permission: a resource with its granted
// scopes and arbitrary claims.
type Permission struct {
	Resource string      // maps to the token's rsname
	Scopes   []string
	Claims   interface{} // e.g. struct{ AppSlug, BuildID []string }
}

// TokenOptions configures a minted token. Audience is required; Issuer defaults
// to the server's Issuer() and Expiry to 24h from now.
type TokenOptions struct {
	Audience    []string
	Issuer      string
	Expiry      time.Time
	Scopes      []string // top-level space-joined "scope" claim
	Permissions []Permission
}

// The JSON shape below mirrors the umaToken structs the service validator
// unmarshals. The round-trip test in this package validates a minted token
// through the real validator, so any drift in the validator's expected shape
// breaks CI here.
type umaClaims struct {
	Issuer        string           `json:"iss,omitempty"`
	Audience      jwt.Audience     `json:"aud,omitempty"`
	IssuedAt      *jwt.NumericDate `json:"iat,omitempty"`
	Expiry        *jwt.NumericDate `json:"exp,omitempty"`
	Scope         string           `json:"scope,omitempty"`
	Authorization *umaAuthz        `json:"authorization,omitempty"`
}

type umaAuthz struct {
	Permissions []umaPermission `json:"permissions,omitempty"`
}

type umaPermission struct {
	Scopes []string    `json:"scopes,omitempty"`
	Claims interface{} `json:"claims,omitempty"`
	Rsname string      `json:"rsname,omitempty"`
}

// Token mints a UMA JWT signed with the server's private key.
func (s *Server) Token(o TokenOptions) (string, error) {
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: s.key},
		(&jose.SignerOptions{ExtraHeaders: map[jose.HeaderKey]interface{}{"kid": s.kid}}).WithType("JWT"),
	)
	if err != nil {
		return "", fmt.Errorf("create signer: %w", err)
	}

	issuer := o.Issuer
	if issuer == "" {
		issuer = s.Issuer()
	}
	expiry := o.Expiry
	if expiry.IsZero() {
		expiry = time.Now().Add(24 * time.Hour)
	}

	claims := umaClaims{
		Issuer:   issuer,
		Audience: jwt.Audience(o.Audience),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Expiry:   jwt.NewNumericDate(expiry),
		Scope:    strings.Join(o.Scopes, " "),
	}
	if len(o.Permissions) > 0 {
		authz := &umaAuthz{}
		for _, p := range o.Permissions {
			authz.Permissions = append(authz.Permissions, umaPermission{
				Scopes: p.Scopes,
				Claims: p.Claims,
				Rsname: p.Resource,
			})
		}
		claims.Authorization = authz
	}

	token, err := jwt.Signed(signer).Claims(claims).Serialize()
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return token, nil
}
