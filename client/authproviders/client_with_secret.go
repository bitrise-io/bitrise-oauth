package authproviders

import (
	"context"
	"net/http"

	"github.com/bitrise-io/bitriseoauth/client"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// ClientWithSecret is a *http.Client preconfigured with Client ID and Client Secret based Oauth2.0 authentication.
// TokenURL is exported to make it it for development/debugging purposes.
type ClientWithSecret struct {
	clientID     string
	clientSecret string
	TokenURL     string
}

// NewWithSecret will return the preconfigured model.
func NewWithSecret(clientID, clientSecret string) client.AuthProvider {
	return ClientWithSecret{
		TokenURL:     "http://104.154.234.133/auth/realms/master/protocol/openid-connect/token",
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// Client is a preconfigured http client using Background context.
func (kcs ClientWithSecret) Client() *http.Client {
	creds := clientcredentials.Config{
		ClientID:     kcs.clientID,
		ClientSecret: kcs.clientSecret,
		TokenURL:     kcs.TokenURL,
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	return creds.Client(context.Background())
}
