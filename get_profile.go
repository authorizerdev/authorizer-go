package authorizer

import (
	"encoding/json"
	"fmt"
)

// GetProfile is method attached to AuthorizerClient.
// It performs profile query on authorizer instance.
// It returns User reference or error.
// For implementation details check GetProfileExample examples/get_profile.go
func (c *AuthorizerClient) GetProfile(headers map[string]string) (*User, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`query {	profile { %s } }`, UserFragment),
		Variables: nil,
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]*User
	json.Unmarshal(bytesData, &res)

	return res["profile"], nil
}
