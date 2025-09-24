package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// JwtValidatorRepository contains a set of JWT validators and can return the appropriate one for a given request
// The request must contain a valid JWT in the Authorization header  ("Authorization: Bearer <token>")
// The validator is selected based on the "iss" claim in the JWT
type JwtValidatorRepository interface {
	GetJwtValidatorForRequest(r *http.Request) (Validator, error)
}

// DefaultJwtValidatorRepository ...
type DefaultJwtValidatorRepository struct {
	JwtValidators map[string]Validator
}

// NewJwtValidatorRepository Creates a new JwtValidatorContainer that holds a set of validators associated with their issuer (iss)
func NewJwtValidatorRepository(jwtValidators map[string]Validator) JwtValidatorRepository {
	return &DefaultJwtValidatorRepository{
		JwtValidators: jwtValidators,
	}
}

// GetJwtValidatorForRequest ...
func (vr *DefaultJwtValidatorRepository) GetJwtValidatorForRequest(r *http.Request) (Validator, error) {
	rawJwt := strings.Split(strings.TrimSpace(r.Header.Get("Authorization")), "Bearer ")
	if len(rawJwt) != 2 {
		return nil, errors.New("failed to read JWT from header")
	}

	iss, err := vr.getIssuerFromRawJWT(rawJwt[1])
	if err != nil {
		return nil, errors.Wrap(err, "failed to get issuer form the JWT")
	}

	validator := vr.JwtValidators[iss]
	if validator == nil {
		return nil, errors.New("there is no JWT validator for this issuer")
	}

	return validator, nil
}

func (vr *DefaultJwtValidatorRepository) getIssuerFromRawJWT(rawJwt string) (string, error) {
	jwtParts := strings.Split(rawJwt, ".")

	if len(jwtParts) != 3 {
		return "", errors.New("invalid JWT format")
	}

	data, err := vr.base64Decode(jwtParts[1])
	if err != nil {
		return "", errors.New("failed to decode JWT")
	}

	var payload map[string]interface{}
	err = json.Unmarshal([]byte(data), &payload)
	if err != nil {
		return "", errors.New("failed to unmarshall JWT")
	}

	issuer, ok := payload["iss"].(string)
	if !ok {
		return "", errors.New("there is no issuer in token")
	}

	return issuer, nil
}

func (vr *DefaultJwtValidatorRepository) base64Decode(src string) (string, error) {
	if l := len(src) % 4; l > 0 {
		src += strings.Repeat("=", 4-l)
	}
	decoded, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		return "", fmt.Errorf("decoding error %s", err)
	}
	return string(decoded), nil
}
