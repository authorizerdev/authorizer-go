package authorizer

import (
	"encoding/json"
)

// TODO does not work without cookie

// Logout is method attached to AuthorizerClient.
// It performs Logout mutation on authorizer instance.
// It takes LogoutInput reference as parameter and returns Response reference or error.
// For implementation details check LogoutExample examples/Logout.go
func (c *AuthorizerClient) Logout(headers map[string]string) (*Response, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query:     `mutation { logout { message } }`,
		Variables: nil,
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]*Response
	json.Unmarshal(bytesData, &res)

	return res["logout"], nil
}
