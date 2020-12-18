package client

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bitrise-io/bitrise-oauth/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	umaGrantType        = "urn:ietf:params:oauth:grant-type:uma-ticket"
	umaClaimTokenFormat = "urn:ietf:params:oauth:token-type:jwt"

	grantType        = "grant_type"
	claimToken       = "claim_token"
	claimTokenFormat = "claim_token_format"
	clientID         = "client_id"
	clientSecret     = "client_secret"
	permission       = "permission"
	audience         = "audience"

	contentType    = "Content-Type"
	formURLEncoded = "application/x-www-form-urlencoded"
)

type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
}

func (e *tokenJSON) expiry(baseTime time.Time) time.Time {
	if v := e.ExpiresIn; v != 0 {
		return baseTime.Add(time.Duration(v) * time.Second)
	}

	return time.Now()
}

// UMATokenSource represents an UMA token source.
type UMATokenSource interface {
	Token(claim interface{}, permisson []Permission, audienceConfig config.AudienceConfig) (*oauth2.Token, error)
}

type umaTokenSource struct {
	config  clientcredentials.Config
	options []UMATokenSourceOption
}

// NewUMATokenSource returns a new UMA token source.
func newUMATokenSource(config clientcredentials.Config, options ...UMATokenSourceOption) umaTokenSource {
	return umaTokenSource{
		config:  config,
		options: options,
	}
}

// Token returns a new UMA token upon each invocation.
func (tokenSource umaTokenSource) Token(claim interface{}, permisson []Permission, audienceConfig config.AudienceConfig) (*oauth2.Token, error) {
	tokenSourceConfig := &UMATokenSourceConfig{}
	for _, opt := range tokenSource.options {
		opt(tokenSourceConfig)
	}

	encodedClaim, err := encodeClaim(claim)
	if err != nil {
		return nil, err
	}

	if tokenSourceConfig.audience != nil {
		audienceConfig = mergeAudienceConfigs(*tokenSourceConfig.audience, audienceConfig)
	}

	request, err := tokenSource.newTokenRequest(encodedClaim, permisson, audienceConfig)
	if err != nil {
		return nil, err
	}

	response, err := sendRequest(request)
	if err != nil {
		return nil, err
	}

	body, err := extractResponseBody(response)
	if err != nil {
		return nil, err
	}

	token, err := extractTokenFromBody(body, time.Now())
	if err != nil {
		return nil, err
	}

	return token, nil
}

func mergeAudienceConfigs(a, b config.AudienceConfig) config.AudienceConfig {
	if len(b.All()) == 0 {
		b = a
	} else {
		allAudiences := append(b.All(), a.All()...)
		if len(allAudiences) > 1 {
			b = config.NewAudienceConfig(allAudiences[0], allAudiences[1:]...)
		} else {
			b = config.NewAudienceConfig(allAudiences[0])
		}
	}
	return b
}

func encodeClaim(claim interface{}) (string, error) {
	bytes, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}

	return b64.StdEncoding.EncodeToString(bytes), nil
}

func (tokenSource umaTokenSource) newTokenRequest(encodedClaim string, permisson []Permission, audienceConfig config.AudienceConfig) (*http.Request, error) {
	v := url.Values{}

	v.Set(grantType, umaGrantType)
	v.Set(claimToken, encodedClaim)
	v.Set(claimTokenFormat, umaClaimTokenFormat)
	v.Set(clientID, tokenSource.config.ClientID)
	v.Set(clientSecret, tokenSource.config.ClientSecret)

	for _, p := range permisson {
		v.Add(permission, p.requestParam())
	}

	for _, a := range audienceConfig.All() {
		v.Add(audience, a)
	}

	request, err := http.NewRequest(http.MethodPost, tokenSource.config.TokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set(contentType, formURLEncoded)

	return request, nil
}

func sendRequest(request *http.Request) (*http.Response, error) {
	client := http.Client{}
	return client.Do(request)
}

func extractResponseBody(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}

	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, &oauth2.RetrieveError{
			Response: response,
			Body:     body,
		}
	}

	return body, nil
}

func extractTokenFromBody(body []byte, baseTime time.Time) (*oauth2.Token, error) {
	var tj tokenJSON

	if err := json.Unmarshal(body, &tj); err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  tj.AccessToken,
		TokenType:    tj.TokenType,
		RefreshToken: tj.RefreshToken,
		Expiry:       tj.expiry(baseTime),
	}

	return token, nil
}
