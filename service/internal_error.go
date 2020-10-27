package service

import (
	"fmt"
	"strings"
)

const (
	// BitriseOAuthInternalError is used to identify the package's internal errors
	BitriseOAuthInternalError = "Bitrise OAuth internal error"
)

// InternalError ...
type InternalError struct {
	Err error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("%s: %v", BitriseOAuthInternalError, e.Err)
}

// IsInternalError returns TRUE if the error is an internal error
func IsInternalError(err error) bool {
	return strings.Contains(err.Error(), BitriseOAuthInternalError)
}
