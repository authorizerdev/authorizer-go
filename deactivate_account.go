package authorizer

import (
	"encoding/json"
)

// DeactivateAccount is method attached to AuthorizerClient.
// It performs deactivate_account mutation on authorizer instance.
// It returns Response reference or error.
// For implementation details check DeactivateAccountExample examples/deactivate_account.go
func (c *AuthorizerClient) DeactivateAccount(headers map[string]string) (*Response, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query:     `mutation deactivateAccount { deactivate_account { message } }`,
		Variables: nil,
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]*Response
	json.Unmarshal(bytesData, &res)

	return res["deactivate_account"], nil
}
