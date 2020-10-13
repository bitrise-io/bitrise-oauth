package authproviders_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bitrise-io/bitrise-oauth/client/authproviders"
	"github.com/bitrise-io/bitrise-oauth/mocks"
	"github.com/stretchr/testify/assert"
)

func Example() {
	authproviders.NewClientWithSecret("my_client_id", "my_client_secret").Client()
}

type tokenJSON struct {
	AccessToken  string        `json:"access_token"`
	TokenType    string        `json:"token_type"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    time.Duration `json:"expires_in"` // at least PayPal returns string, while most return number
}

func TestNewClientWithSecret_threads_using_same_client(t *testing.T) {
	clientsToCreate := 30
	callsPerClient := 30

	var createdClients sync.Map

	var wg sync.WaitGroup
	wg.Add(clientsToCreate * callsPerClient)

	for i := 0; i < callsPerClient; i++ {
		go func() {
			for j := 0; j < clientsToCreate; j++ {
				go func(j int) {
					defer wg.Done()
					c := authproviders.NewClientWithSecret(fmt.Sprintf("clientID-%d", j), fmt.Sprintf("clientSecret-%d", j),
						authproviders.WithTokenURL("myurl")).Client()

					pointerAddress := fmt.Sprintf("%p", c)
					createdClients.Store(pointerAddress, c)
				}(j)
			}
		}()
	}

	wg.Wait()

	createdClientsCount := 0
	createdClients.Range(func(_, _ interface{}) bool {
		createdClientsCount++
		return true
	})

	assert.Equal(t, clientsToCreate, createdClientsCount)
}

func TestNewClientWithSecret_not_using_refresh_token(t *testing.T) {
	mockedAuthService := mocks.AuthService{}
	mockedClient := mocks.Client{}

	accessToken, refreshToken := "initial-access-token", "initial-refresh-token"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Add("content-type", "application/json")

			assert.NoError(t, r.ParseForm())
			grantType := r.PostForm.Get("grant_type")

			if grantType == "refresh_token" {
				accessToken, refreshToken = "refreshed-access-token", "refreshed-refresh-token"
			}

			assert.NoError(t, json.NewEncoder(w).Encode(tokenJSON{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				TokenType:    "Bearer",
				ExpiresIn:    11, // go has a -10 seconds delta time gap - https://github.com/golang/oauth2/blob/master/token.go#L22
			}))

			mockedAuthService.Token(grantType)
		default:
			tokenHeaderSplit := strings.Split(r.Header.Get("Authorization"), " ")
			assert.Len(t, tokenHeaderSplit, 2)

			mockedClient.Test(tokenHeaderSplit[1])

			w.WriteHeader(http.StatusOK)
		}
	}))
	defer ts.Close()

	mockedAuthService.
		On("Token", "client_credentials").Return().
		Twice()
	mockedClient.
		On("Test", "initial-access-token").Return().
		Times(6)

	c := authproviders.NewClientWithSecret("my-client-id", "my-secret",
		authproviders.WithTokenURL(ts.URL+"/token")).Client()

	for i := 0; i < 6; i++ {
		resp, err := c.Get(ts.URL + "/test")
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		time.Sleep(time.Millisecond * 400)
	}

	mockedAuthService.AssertExpectations(t)
	mockedClient.AssertExpectations(t)
	mockedAuthService.AssertNotCalled(t, "Token", "refresh_token")
	mockedClient.AssertNotCalled(t, "Test", "refreshed-access-token")
}
