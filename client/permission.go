package client

import "fmt"

// Permission ...
type Permission struct {
	resourceName       string
	authorizationScope string
}

// NewPermission ...
func NewPermission(resourceName, authorizationScope string) Permission {
	return Permission{
		resourceName:       resourceName,
		authorizationScope: authorizationScope,
	}
}

func (permission Permission) requestParam() string {
	return fmt.Sprintf("%s#%s", permission.resourceName, permission.authorizationScope)
}
