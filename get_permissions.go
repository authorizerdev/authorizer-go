package authorizer

import (
	"encoding/json"
	"fmt"
)

// GetPermissions is method attached to AuthorizerClient.
// It performs the permissions query on the authorizer instance and returns
// the list of resource:scope permissions granted to the authenticated
// principal. The principal is identified from the request headers
// (Authorization bearer token or session cookie), hence headers are required.
// It returns a slice of Permission or error.
func (c *AuthorizerClient) GetPermissions(headers map[string]string) ([]Permission, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`query permissions { permissions { %s } }`, PermissionFragment),
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string][]Permission
	json.Unmarshal(bytesData, &res)

	return res["permissions"], nil
}
