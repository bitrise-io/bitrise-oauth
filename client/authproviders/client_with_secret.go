package authproviders

import (
	"context"
	"net/http"

	"github.com/bitrise-io/bitrise-oauth/client"
	"github.com/bitrise-io/bitrise-oauth/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// ClientWithSecret is a *http.Client preconfigured with Client ID and Client Secret based Oauth2.0 authentication.
// TokenURL is exported to make it it for development/debugging purposes.
type ClientWithSecret struct {
	clientID     string
	clientSecret string
	tokenURL     string
}

// NewClientWithSecret will return the preconfigured model.
func NewClientWithSecret(clientID, clientSecret string, opts ...ClientOption) client.AuthProvider {
	cws := &ClientWithSecret{
		tokenURL:     config.TokenURL,
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	for _, opt := range opts {
		opt(cws)
	}

	return cws
}

// Client is a preconfigured http client using Background context.
func (kcs ClientWithSecret) Client() *http.Client {
	creds := clientcredentials.Config{
		ClientID:     kcs.clientID,
		ClientSecret: kcs.clientSecret,
		TokenURL:     kcs.tokenURL,
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	return creds.Client(context.Background())
}
