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
	umaGrantType     = "urn:ietf:params:oauth:grant-type:uma-ticket"
	claimTokenFormat = "urn:ietf:params:oauth:token-type:jwt"
)

type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
}

func (e *tokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	return
}

// UMATokenSource returns a token source that returns a new token each time the Token()
// method is called.
type UMATokenSource interface {
	Token(claim interface{}, permisson []Permission, audienceConfig config.AudienceConfig) (*oauth2.Token, error)
}

type umaTokenSource struct {
	config clientcredentials.Config
}

// NewUMATokenSource returns a new UMA token source.
func NewUMATokenSource(config clientcredentials.Config) UMATokenSource {
	return umaTokenSource{
		config: config,
	}
}

// Token returns a new UMA token upon each invocation.
func (tokenSource umaTokenSource) Token(claim interface{}, permisson []Permission, audienceConfig config.AudienceConfig) (*oauth2.Token, error) {
	encodedClaim, err := encodeClaim(claim)
	if err != nil {
		return nil, err
	}

	request, err := newTokenRequest(tokenSource.config, encodedClaim, permisson, audienceConfig)
	if err != nil {
		return nil, err
	}

	token, err := sendRequest(request)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func encodeClaim(claim interface{}) (string, error) {
	bytes, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}

	return b64.StdEncoding.EncodeToString(bytes), nil
}

func newTokenRequest(config clientcredentials.Config, encodedClaim string, permisson []Permission, audienceConfig config.AudienceConfig) (*http.Request, error) {
	v := url.Values{}

	v.Set("grant_type", umaGrantType)
	v.Set("claim_token", encodedClaim)
	v.Set("claim_token_format", claimTokenFormat)
	v.Set("client_id", config.ClientID)
	v.Set("client_secret", config.ClientSecret)

	for _, p := range permisson {
		v.Set("permission", p.requestParam())
	}

	for _, a := range audienceConfig.All() {
		v.Set("audience", a)
	}

	request, err := http.NewRequest("POST", config.TokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}

func sendRequest(request *http.Request) (*oauth2.Token, error) {
	client := http.Client{}
	resp, err := client.Do(request)

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}

	if code := resp.StatusCode; code < 200 || code > 299 {
		return nil, &oauth2.RetrieveError{
			Response: resp,
			Body:     body,
		}
	}

	var tj tokenJSON
	if err = json.Unmarshal(body, &tj); err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  tj.AccessToken,
		TokenType:    tj.TokenType,
		RefreshToken: tj.RefreshToken,
		Expiry:       tj.expiry(),
	}

	return token, nil
}
