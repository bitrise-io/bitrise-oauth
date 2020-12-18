package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/bitrise-io/bitrise-oauth/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// AuthProvider which creates an authenticated *http.Client ready to use with our authentication provider.
type AuthProvider interface {
	ManagedHTTPClient(...HTTPClientOption) *http.Client
	HTTPClient(...HTTPClientOption) *http.Client
	TokenSource() oauth2.TokenSource
	UMATokenSource(options ...UMATokenSourceOption) UMATokenSource
}

var clients sync.Map

// WithSecret hold the ouath client credentials
type WithSecret struct {
	clientID     string
	clientSecret string
	realm        string
	baseURL      string
	credentials  clientcredentials.Config
	scopes       []string
}

// NewWithSecret will return the preconfigured model.
func NewWithSecret(clientID, clientSecret string, scopeOption ScopeOption, opts ...Option) AuthProvider {
	cws := &WithSecret{
		baseURL:      config.BaseURL,
		realm:        config.Realm,
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	for _, opt := range opts {
		opt(cws)
	}

	scopeOption(cws)

	cws.credentials = clientcredentials.Config{
		ClientID:     cws.clientID,
		ClientSecret: cws.clientSecret,
		TokenURL:     cws.tokenURL(),
		Scopes:       cws.scopes,
	}

	return cws
}

func (cws *WithSecret) tokenURL() string {
	return fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/token", cws.baseURL, cws.realm)
}

func (cws *WithSecret) uid() string {
	return strings.Join([]string{cws.clientID, cws.clientSecret, cws.tokenURL()}, "-")
}

// TokenSource returns a token source that refreshes the token only when expires
func (cws *WithSecret) TokenSource() oauth2.TokenSource {
	return cws.credentials.TokenSource(context.Background())
}

// UMATokenSource returns an UMA token source.
func (cws *WithSecret) UMATokenSource(options ...UMATokenSourceOption) UMATokenSource {
	return newUMATokenSource(cws.credentials, options...)
}

// ManagedHTTPClient is a preconfigured http client using in-memory client storage
// this way the clients with the same credentials will be reused.
func (cws *WithSecret) ManagedHTTPClient(opts ...HTTPClientOption) *http.Client {
	client := cws.HTTPClient(opts...)
	storedClient, _ := clients.LoadOrStore(cws.uid(), client)
	return storedClient.(*http.Client)
}

// HTTPClient is a preconfigured http client
func (cws *WithSecret) HTTPClient(opts ...HTTPClientOption) *http.Client {
	creds := clientcredentials.Config{
		ClientID:     cws.clientID,
		ClientSecret: cws.clientSecret,
		TokenURL:     cws.tokenURL(),
	}

	clientOpts := &HTTPClientConfig{
		context: context.Background(),
	}

	for _, opt := range opts {
		opt(clientOpts)
	}

	client := &http.Client{}
	if clientOpts.baseClient != nil {
		client = clientOpts.baseClient
	}

	origTransport := client.Transport

	resettableTokenSrc := &resettableTokenSource{
		ctx:   clientOpts.context,
		src:   creds.TokenSource(clientOpts.context),
		creds: creds}

	invalidTokenRefresherTransport := &invalidTokenRefresherTransport{
		base: &oauth2.Transport{
			Source: resettableTokenSrc,
			Base:   origTransport},
		tokenSrc: resettableTokenSrc,
	}

	client.Transport = invalidTokenRefresherTransport

	return client
}
