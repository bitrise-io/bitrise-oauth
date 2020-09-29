package authproviders

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// KeycloakClientWithSecret ...
type KeycloakClientWithSecret struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
}

// Client ...
func (kcs KeycloakClientWithSecret) Client() *http.Client {
	creds := clientcredentials.Config{
		ClientID:     kcs.ClientID,
		ClientSecret: kcs.ClientSecret,
		TokenURL:     kcs.TokenURL,
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	return creds.Client(context.Background())
}
