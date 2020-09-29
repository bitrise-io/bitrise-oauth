package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	gocloak "github.com/Nerzal/gocloak/v7"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func testCall() {
	time.Sleep(time.Second * 3)
	fmt.Println("# get")
	creds := clientcredentials.Config{
		ClientID:     "tomi-test",
		ClientSecret: "37dd1bc5-50bb-4674-a5fa-2cec87037e52",
		TokenURL:     "http://104.154.234.133/auth/realms/master/protocol/openid-connect/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	resp, err := creds.Client(context.Background()).Get("http://localhost:8080/test")
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

// KeycloakTokenIntrospector ...
type KeycloakTokenIntrospector struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

// Middleware ...
func (kti KeycloakTokenIntrospector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authKey := r.Header.Get("Authorization")
		if s := strings.Split(authKey, " "); len(s) == 2 {
			authKey = s[1]
		}
		c := gocloak.NewClient(kti.BaseURL)
		result, err := c.RetrospectToken(context.Background(), authKey, kti.ClientID, kti.ClientSecret, kti.Realm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if result == nil || result.Active == nil || !*result.Active {
			b, _ := json.MarshalIndent(result, "", " ")
			http.Error(w, string(b), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	go testCall()

	fmt.Println("listening on :3333")
	mux := http.NewServeMux()

	kti := KeycloakTokenIntrospector{
		BaseURL:      "http://104.154.234.133/",
		Realm:        "master",
		ClientID:     "tomi-test-backend",
		ClientSecret: "a7d065c5-0a48-4b42-922e-82636bba85fe",
	}

	mux.Handle("/test", kti.Middleware(http.HandlerFunc(handler)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
