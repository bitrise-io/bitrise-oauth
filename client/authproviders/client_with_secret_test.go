package authproviders_test

import "github.com/bitrise-io/bitriseoauth/client/authproviders"

func Example() {
	authproviders.NewWithSecret("my_client_id", "my_client_secret").Client()
}
