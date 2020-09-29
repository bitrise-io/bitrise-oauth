package authprovider

import (
	"context"
	"net/http"

	"github.com/bitrise-io/bitriseoauth/settings"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const baseURL = `http://localhost:4444`

// Hydra ...
type Hydra struct{}

// ClientWithIDAndSecret ...
func (h Hydra) ClientWithIDAndSecret(id, secret string) *http.Client {
	creds := clientcredentials.Config{
		ClientID:     "my-client",
		ClientSecret: "secret",
		TokenURL:     "http://localhost:4444/oauth2/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	return creds.Client(context.Background())
}

func (h Hydra) ClientSettings() settings.ClientSettings {
	return h
}
