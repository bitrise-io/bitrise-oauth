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
}

var clients sync.Map

// WithSecret hold the ouath client credentials
type WithSecret struct {
	clientID     string
	clientSecret string
	realm        string
	baseURL      string
	credentials  clientcredentials.Config
}

// NewWithSecret will return the preconfigured model.
func NewWithSecret(clientID, clientSecret string, opts ...Option) AuthProvider {
	cws := &WithSecret{
		baseURL:      config.BaseURL,
		realm:        config.Realm,
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	for _, opt := range opts {
		opt(cws)
	}

	cws.credentials = clientcredentials.Config{
		ClientID:     cws.clientID,
		ClientSecret: cws.clientSecret,
		TokenURL:     cws.tokenURL(),
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

	if clientOpts.baseClient == nil {
		return creds.Client(clientOpts.context)
	}

	oauth2Transport := &oauth2.Transport{
		Source: creds.TokenSource(clientOpts.context),
		Base:   clientOpts.baseClient.Transport,
	}

	clientOpts.baseClient.Transport = oauth2Transport

	return clientOpts.baseClient
}
