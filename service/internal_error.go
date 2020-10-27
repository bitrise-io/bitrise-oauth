package service

import (
	"fmt"
	"strings"
)

const (
	bitriseOAuthInternalError = "Bitrise OAuth internal error"
)

// InternalError ...
type InternalError struct {
	Err error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("%s: %v", bitriseOAuthInternalError, e.Err)
}

// IsInternalError returns TRUE if the error is an internal error
func IsInternalError(err error) bool {
	return strings.Contains(err.Error(), bitriseOAuthInternalError)
}
