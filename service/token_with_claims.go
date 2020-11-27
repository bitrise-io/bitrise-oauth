package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	authorizationKey = "authorization"
	permissionsKey   = "permissions"
)

type umaToken struct {
	Issuer        string          `json:"iss,omitempty"`
	Audience      jwt.Audience    `json:"aud,omitempty"`
	IssuedAt      jwt.NumericDate `json:"iat,omitempty"`
	Expiry        jwt.NumericDate `json:"exp,omitempty"`
	Authorization authorization   `json:"authorization,omitempty"`
}

type authorization struct {
	Permissions []permisson `json:"permissions,omitempty"`
}

type permisson struct {
	Scopes []string    `json:"scopes,omitempty"`
	Claims interface{} `json:"claims,omitempty"`
	Rsid   string      `json:"rsid,omitempty"`
	Rsname string      `json:"rsname,omitempty"`
}

// TokenWithClaims is a wrapper over jwt.JSONWebToken to extract the
// claims easily.
type TokenWithClaims struct {
	key   interface{}
	token *jwt.JSONWebToken
}

// Payload returns the  contents of the token.
func (tokenWithClaim *TokenWithClaims) Payload() (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	if err := tokenWithClaim.token.Claims(tokenWithClaim.key, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}

// Permissions returns the persmissions part of the token.
func (tokenWithClaim *TokenWithClaims) Permissions() ([]interface{}, error) {
	payload, err := tokenWithClaim.Payload()
	if err != nil {
		return nil, err
	}

	authorization, ok := payload[authorizationKey].(map[string]interface{})
	if !ok {
		return nil, errors.New("authorization is missing from token")
	}

	permissions, ok := authorization[permissionsKey].([]interface{})
	if !ok {
		return nil, errors.New("permissions is missing from token")
	}

	return permissions, nil
}

// Claim returns the claim for the provided resource's name.
func (tokenWithClaim *TokenWithClaims) Claim(resourceName string, claim interface{}) error {
	token := umaToken{}
	if err := tokenWithClaim.token.Claims(tokenWithClaim.key, &token); err != nil {
		return err
	}

	for _, permission := range token.Authorization.Permissions {
		if permission.Rsname == resourceName {
			// First we have to serialize to json
			jsonClaims, err := json.Marshal(permission.Claims)
			if err != nil {
				return err
			}

			// Deserialize from json
			if err := json.Unmarshal(jsonClaims, &claim); err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("permission for resource: %s not found", resourceName)
}
