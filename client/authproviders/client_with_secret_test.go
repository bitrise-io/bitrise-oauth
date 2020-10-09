package authproviders_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"
	"time"

	"github.com/bitrise-io/bitrise-oauth/client/authproviders"
	"github.com/c2fo/testify/mock"
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

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Token(grantType string) {
	m.Called(grantType)
}

type mockService struct {
	mock.Mock
}

func (m *mockService) Test(accessToken string) {
	m.Called(accessToken)
}

func TestNewClientWithSecret_threads_not_using_same_client(t *testing.T) {
	mockedAuthService := mockAuthService{}
	mockedService := mockService{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Add("content-type", "application/json")

			assert.NoError(t, json.NewEncoder(w).Encode(tokenJSON{
				AccessToken:  "my-access-token",
				RefreshToken: "my-refresh-token",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			}))

			assert.NoError(t, r.ParseForm())
			mockedAuthService.Token(r.PostForm.Get("grant_type"))
		default:
			mockedService.Test("")
			w.WriteHeader(http.StatusOK)
		}

	}))
	defer ts.Close()

	mockedAuthService.
		On("Token", "client_credentials").Return().
		Once()
	mockedService.
		On("Test", "").Return().Times(5)

	for i := 0; i < 5; i++ {
		// TODO: remove this comment -> move the two lines above to for loop so the test will pass
		c := authproviders.NewClientWithSecret("my-client-id", "my-secret",
			authproviders.WithCustomTokenURL(ts.URL+"/token")).Client()

		resp, err := c.Get(ts.URL + "/test")
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusOK)
	}

	mockedAuthService.AssertExpectations(t)
	mockedService.AssertExpectations(t)
}

func TestNewClientWithSecret_using_refresh_token(t *testing.T) {
	mockedAuthService := mockAuthService{}
	mockedService := mockService{}

	accessToken, refreshToken := "initial-access-token", "initial-refresh-token"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := httputil.DumpRequest(r, true)
		fmt.Println(string(b))
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

			mockedService.Test(tokenHeaderSplit[1])

			w.WriteHeader(http.StatusOK)
		}
	}))
	defer ts.Close()

	mockedAuthService.
		On("Token", "client_credentials").Return().
		Twice()
	mockedService.
		On("Test", "initial-access-token").Return().
		Times(6)

	c := authproviders.NewClientWithSecret("my-client-id", "my-secret",
		authproviders.WithCustomTokenURL(ts.URL+"/token")).Client()

	for i := 0; i < 6; i++ {
		resp, err := c.Get(ts.URL + "/test")
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		time.Sleep(time.Millisecond * 400)
	}

	mockedAuthService.AssertExpectations(t)
	mockedService.AssertExpectations(t)
	mockedAuthService.AssertNotCalled(t, "Token", "refresh_token")
	mockedService.AssertNotCalled(t, "Test", "refreshed-access-token")
}
