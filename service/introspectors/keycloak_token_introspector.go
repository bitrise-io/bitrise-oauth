package introspectors

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	gocloak "github.com/Nerzal/gocloak/v7"
)

// KeycloakToken ...
type KeycloakToken struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

// Middleware ...
func (kti KeycloakToken) Middleware(next http.Handler) http.Handler {
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
