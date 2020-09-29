package client_test

import (
	"testing"

	"github.com/bitrise-io/bitriseoauth/client"
	"github.com/bitrise-io/bitriseoauth/options"
)

func TestClient(t *testing.T) {
	_ = client.New(options.WithClientIDAndSecret{
		ClientID: "", ClientSecret: "",
	})

}
