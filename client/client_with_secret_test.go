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

	"github.com/bitrise-io/bitrise-oauth/config"

	"github.com/bitrise-io/bitrise-oauth/client"
	"github.com/bitrise-io/bitrise-oauth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Example() {
	client.NewWithSecret("my_client_id", "my_client_secret", client.WithScope("")).ManagedHTTPClient()
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
		c := client.NewWithSecret(fmt.Sprintf("clientID-%d", i), fmt.Sprintf("clientSecret-%d", i), client.WithScope("scope"),
			client.WithBaseURL("https://google.com"), client.WithRealm("myrealm")).ManagedHTTPClient()

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

func Test_GivenDifferentClientConfigs_WhenTheManagedHttpClientsAreInstantiated_ThenExpectNewClientsToBeCreated(t *testing.T) {
	// Given
	configs := []struct {
		clientID     string
		clientSecret string
		scope        string
		realm        string
		baseURL      string
	}{
		{"id1", "secret1", "scope1", "realm1", "https://url1.com"},
		{"id1", "secret1", "scope2", "realm1", "https://url1.com"},
		{"id1", "secret1", "scope1", "realm2", "https://url1.com"},
		{"id1", "secret1", "scope1", "realm1", "https://url2.com"},
	}

	// When
	var createdClients []*http.Client
	for _, conf := range configs {
		c := client.NewWithSecret(conf.clientID,
			conf.clientSecret,
			client.WithScope(conf.scope),
			client.WithBaseURL(conf.baseURL),
			client.WithRealm(conf.realm)).ManagedHTTPClient()
		createdClients = append(createdClients, c)
	}

	// Then
	for i := 0; i < len(createdClients); i++ {
		for j := i + 1; j < len(createdClients); j++ {
			assert.NotEqual(t, createdClients[i], createdClients[j])
		}
	}
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

	ts := startMockServer(t, &mockedAuthService, &mockedClient, accessToken, http.StatusOK, http.StatusOK)
	defer ts.Close()

	mockedAuthService.
		On("Token").Return().
		Twice()
	mockedClient.
		On("Test", accessToken+"-0").Return().
		Times(3)
	mockedClient.
		On("Test", accessToken+"-1").Return().
		Times(3)

	// When
	c := client.NewWithSecret("my-client-id", "my-secret", client.WithScope(""),
		client.WithBaseURL(ts.URL)).ManagedHTTPClient()

	// Then
	for i := 0; i < 6; i++ {
		resp, err := c.Get(ts.URL)
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		time.Sleep(time.Millisecond * 400)
	}

	mockedAuthService.AssertExpectations(t)
	mockedClient.AssertExpectations(t)
}

func Test_GivenAnExistingHTTPClient_WhenItIsPassedAsAnOptionDuringInstantiation_ThenExpectTheNewClientToBeAnExtendedCopyOfTheExistingOne(t *testing.T) {
	// Given
	baseClient := &http.Client{}

	// When
	client := client.NewWithSecret("test-id", "test-secret", client.WithScope("")).HTTPClient(client.WithBaseClient(baseClient))

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
	client := client.NewWithSecret("test-id", "test-secret", client.WithScope(""), client.WithBaseURL(ts.URL)).HTTPClient(client.WithContext(baseCtx))

	url := ts.URL

	// Then
	_, err := client.Get(url)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf(`Get "%s": context canceled`, url))
}

func startMockServer(t *testing.T, mockedAuthService *mocks.AuthService, mockedClient *mocks.Client, accessToken string, tokenStatusCode, defaultStatusCode int) *httptest.Server {
	tokenEndpointURL := "/auth/realms/" + config.Realm + "/protocol/openid-connect/token"

	counter := 0
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case tokenEndpointURL:
			w.Header().Add("content-type", "application/json")

			assert.NoError(t, json.NewEncoder(w).Encode(tokenJSON{
				AccessToken:  fmt.Sprintf("%s-%d", accessToken, counter),
				RefreshToken: "refresh-token",
				TokenType:    "Bearer",
				ExpiresIn:    11, // go has a -10 seconds delta time gap - https://github.com/golang/oauth2/blob/master/token.go#L22
			}))

			counter++

			mockedAuthService.Token()

			w.WriteHeader(tokenStatusCode)
		default:
			tokenHeaderSplit := strings.Split(r.Header.Get("Authorization"), " ")
			assert.Len(t, tokenHeaderSplit, 2)

			mockedClient.Test(tokenHeaderSplit[1])

			w.WriteHeader(defaultStatusCode)
		}
	}))
}

func Test_GivenTokenSourceWithTokenThatWillNotExpireBetweenRequests_WhenTokenStoreIsFetchedMultipleTimes_ThenExpectTheSameTokenGranted(t *testing.T) {
	// Given
	mockedAuthService := &mocks.AuthService{}
	mockedClient := &mocks.Client{}

	accessToken := "initial-access-token"

	ts := startMockServer(t, mockedAuthService, mockedClient, accessToken, http.StatusOK, http.StatusOK)
	defer ts.Close()

	mockedAuthService.
		On("Token").Return().
		Once()

	tokenSource := client.NewWithSecret("my-client-id", "my-secret", client.WithScope(""),
		client.WithBaseURL(ts.URL)).TokenSource()

	// When
	token, err := tokenSource.Token()
	require.NoError(t, err)
	require.Equal(t, token.AccessToken, accessToken+"-0")

	token, err = tokenSource.Token()
	require.NoError(t, err)
	require.Equal(t, token.AccessToken, accessToken+"-0")

	// Then
	mockedAuthService.AssertExpectations(t)
	mockedClient.AssertExpectations(t)
}

func Test_GivenAServerThatRejectsHTTPCall_WhenAGetCallIsFired_ThenExpectTheClientToRenewAccessTokenOnce(t *testing.T) {
	testCases := []struct {
		expectedStatusCode int
		expectedNoOfCalls  int
	}{
		{http.StatusUnauthorized, 2},
		{http.StatusOK, 1},
	}

	for _, testCase := range testCases {
		// Given
		mockedAuthService := mocks.AuthService{}
		mockedClient := mocks.Client{}

		accessToken := "initial-access-token"

		ts := startMockServer(t, &mockedAuthService, &mockedClient, accessToken, http.StatusOK, testCase.expectedStatusCode)
		defer ts.Close()

		mockedAuthService.
			On("Token").Return().
			Times(testCase.expectedNoOfCalls)

		expectAccessTokenChangeForTest(&mockedClient, accessToken, testCase.expectedNoOfCalls)

		// When
		c := client.NewWithSecret("my-client-id", "my-secret", client.WithScope(""),
			client.WithBaseURL(ts.URL)).HTTPClient()

		// Then
		resp, err := c.Get(ts.URL)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)

		mockedAuthService.AssertExpectations(t)
		mockedClient.AssertExpectations(t)
	}
}

func expectAccessTokenChangeForTest(m *mocks.Client, accessTokenPrefix string, expectedNoOfCalls int) {
	for i := 0; i < expectedNoOfCalls; i++ {
		fmt.Println("he", fmt.Sprintf("%s-%d", accessTokenPrefix, i))
		m.On("Test", fmt.Sprintf("%s-%d", accessTokenPrefix, i)).Return().
			Once()
	}
}
