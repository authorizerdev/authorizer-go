package authorizer

import (
	"encoding/json"
	"fmt"
)

// LoginType defines attributes for login request
type LoginInput struct {
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Roles    []*string `json:"roles,omitempty"`
	Scope    []*string `json:"scope,omitempty"`
}

func (c *authorizerClient) Login(req *LoginInput) (*AuthTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`
		mutation login($data: LoginInput!) {
			login(params: $data) {
				%s
			}
		}`, AuthTokenFragment),
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*AuthTokenResponse
	json.Unmarshal(bytesData, &res)

	return res["login"], nil
}
