package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/bitrise-io/bitriseoauth/client"
	"github.com/bitrise-io/bitriseoauth/client/authproviders"
	"github.com/bitrise-io/bitriseoauth/service"
	"github.com/bitrise-io/bitriseoauth/service/introspectors"
)

func testCall() {
	time.Sleep(time.Second * 3)
	fmt.Println("# get")

	var authProvider client.AuthProvider = authproviders.KeycloakClientWithSecret{
		ClientID:     "tomi-test",
		ClientSecret: "37dd1bc5-50bb-4674-a5fa-2cec87037e52",
		TokenURL:     "http://104.154.234.133/auth/realms/master/protocol/openid-connect/token",
	}

	resp, err := authProvider.Client().Get("http://localhost:8080/test")
	if err != nil {
		panic(err)
	}

	rb, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Println("resp:\n" + string(rb))
}

func handler(w http.ResponseWriter, r *http.Request) {
	rb, err := httputil.DumpRequest(r, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rb))
}

func main() {
	go testCall()

	fmt.Println("listening on :3333")
	mux := http.NewServeMux()

	var kti service.Introspector = introspectors.KeycloakToken{
		BaseURL:      "http://104.154.234.133/",
		Realm:        "master",
		ClientID:     "tomi-test-backend",
		ClientSecret: "a7d065c5-0a48-4b42-922e-82636bba85fe",
	}

	mux.Handle("/test", kti.Middleware(http.HandlerFunc(handler)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
