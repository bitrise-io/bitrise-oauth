package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bitrise-io/bitrise-oauth/client"
	"github.com/bitrise-io/bitrise-oauth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Example() {
	client.NewWithSecret("my_client_id", "my_client_secret").ManagedHTTPClient()
}

type tokenJSON struct {
	AccessToken  string        `json:"access_token"`
	TokenType    string        `json:"token_type"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    time.Duration `json:"expires_in"` // at least PayPal returns string, while most return number
}

func Test_Given30ThreadsAndEachWillLaunch30RequestsOnNewThreads_WhenTheManagedHttpClientsAreInstantiated_ThenExpect30HttpClientsToBeCreated(t *testing.T) {
	// Given
	clientsToCreate := 30
	callsPerClient := 20

	var createdClients sync.Map

	// When
	async(clientsToCreate, callsPerClient, func(i, j int) {
		c := client.NewWithSecret(fmt.Sprintf("clientID-%d", i), fmt.Sprintf("clientSecret-%d", i),
			client.WithTokenURL("https://google.com")).ManagedHTTPClient()

		pointerKey := fmt.Sprintf("%d,%d", i, j)
		pointerAddress := fmt.Sprintf("%p", c)
		createdClients.Store(pointerKey, pointerAddress)
	})

	pointerCount := make(map[string]int)

	createdClients.Range(func(k, v interface{}) bool {
		pointerAddress, ok := v.(string)
		if !ok {
			panic("Error in type assertion")
		}
		if entry, found := pointerCount[pointerAddress]; found {
			pointerCount[pointerAddress] = entry + 1
		} else {
			pointerCount[pointerAddress] = 1
		}
		return true
	})

	// Then
	pointerCountLength := len(pointerCount)
	assert.Equal(t, clientsToCreate, pointerCountLength)
	for _, v := range pointerCount {
		assert.Equal(t, callsPerClient, v)
	}
}

func syncMapLen(sm *sync.Map) int {
	len := 0
	sm.Range(func(_, _ interface{}) bool {
		len++
		return true
	})
	return len
}

func async(iCount, jCount int, fn func(int, int)) {
	var wg sync.WaitGroup
	wg.Add(iCount * jCount)
	for i := 0; i < iCount; i++ {
		go func(i int) {
			for j := 0; j < jCount; j++ {
				go func(j int) {
					defer wg.Done()
					fn(i, j)
				}(j)
			}
		}(i)
	}
	wg.Wait()
}

func Test_GivenATokenThatWillExpireAfter1Second_WhenANewTokenIsAcquired_ThenExpectTheRefreshTokenNotToBeUsed(t *testing.T) {
	// Given
	mockedAuthService := mocks.AuthService{}
	mockedClient := mocks.Client{}

	accessToken := "initial-access-token"

	ts := startMockServer(t, &mockedAuthService, &mockedClient, accessToken)
	defer ts.Close()

	mockedAuthService.
		On("Token").Return().
		Twice()
	mockedClient.
		On("Test", "initial-access-token").Return().
		Times(6)

	// When
	c := client.NewWithSecret("my-client-id", "my-secret",
		client.WithTokenURL(ts.URL+"/token")).ManagedHTTPClient()

	// Then
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

func Test_GivenAnExistingHTTPClient_WhenItIsPassedAsAnOptionDuringInstantiation_ThenExpectTheNewClientToBeAnExtendedCopyOfTheExistingOne(t *testing.T) {
	// Given
	baseClient := &http.Client{}

	// When
	client := client.NewWithSecret("test-id", "test-secret").HTTPClient(client.WithBaseClient(baseClient))

	// Then
	assert.Equal(t, baseClient, client)
}

func Test_GivenAnExistingHTTPContext_WhenItIsPassedAsAnOptionDuringInstantiation_ThenExpectTheNewClientToHaveTheSameContextAsTheExistingOne(t *testing.T) {
	// Given
	baseCtx, cancel := context.WithCancel(context.Background())
	cancel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// When
	client := client.NewWithSecret("test-id", "test-secret", client.WithTokenURL(ts.URL+"/token")).HTTPClient(client.WithContext(baseCtx))

	url := ts.URL + "/token"

	// Then
	_, err := client.Get(url)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf(`Get "%s": context canceled`, url))
}

func startMockServer(t *testing.T, mockedAuthService *mocks.AuthService, mockedClient *mocks.Client, accessToken string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Add("content-type", "application/json")

			assert.NoError(t, json.NewEncoder(w).Encode(tokenJSON{
				AccessToken:  accessToken,
				RefreshToken: "refresh-token",
				TokenType:    "Bearer",
				ExpiresIn:    11, // go has a -10 seconds delta time gap - https://github.com/golang/oauth2/blob/master/token.go#L22
			}))

			mockedAuthService.Token()
		default:
			tokenHeaderSplit := strings.Split(r.Header.Get("Authorization"), " ")
			assert.Len(t, tokenHeaderSplit, 2)

			mockedClient.Test(tokenHeaderSplit[1])

			w.WriteHeader(http.StatusOK)
		}
	}))
}

func Test_GivenTokenSourceWithTokenThatWillNotExpireBetweenRequests_WhenTokenStoreIsFetchedMultipleTimes_ThenExpectTheSameTokenGranted(t *testing.T) {
	// Given
	mockedAuthService := &mocks.AuthService{}
	mockedClient := &mocks.Client{}

	accessToken := "initial-access-token"

	ts := startMockServer(t, mockedAuthService, mockedClient, accessToken)
	defer ts.Close()

	mockedAuthService.
		On("Token").Return().
		Once()

	tokenSource := client.NewWithSecret("my-client-id", "my-secret",
		client.WithTokenURL(ts.URL+"/token")).TokenSource()

	// When
	token, err := tokenSource.Token()
	require.NoError(t, err)
	require.Equal(t, token.AccessToken, accessToken)

	token, err = tokenSource.Token()
	require.NoError(t, err)
	require.Equal(t, token.AccessToken, accessToken)

	// Then
	mockedAuthService.AssertExpectations(t)
	mockedClient.AssertExpectations(t)
}
