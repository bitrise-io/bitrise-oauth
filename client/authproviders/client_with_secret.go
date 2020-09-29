package authproviders

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// ClientWithSecret ...
type ClientWithSecret struct {
	clientID     string
	clientSecret string
	TokenURL     string
}

// NewClientWithSecret ...
func NewClientWithSecret(clientID, clientSecret string) ClientWithSecret {
	return ClientWithSecret{
		TokenURL:     "http://104.154.234.133/auth/realms/master/protocol/openid-connect/token",
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// Client ...
func (kcs ClientWithSecret) Client() *http.Client {
	creds := clientcredentials.Config{
		ClientID:     kcs.clientID,
		ClientSecret: kcs.clientSecret,
		TokenURL:     kcs.TokenURL,
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	return creds.Client(context.Background())
}
