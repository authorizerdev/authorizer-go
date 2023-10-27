package authorizer

import (
	"encoding/json"
	"fmt"
)

// LoginInput defines attributes for login request
type LoginInput struct {
	Email       *string   `json:"email,omitempty"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Password    string    `json:"password"`
	Roles       []*string `json:"roles,omitempty"`
	Scope       []*string `json:"scope,omitempty"`
}

// Login is method attached to AuthorizerClient.
// It performs login mutation on authorizer instance.
// It takes LoginInput reference as parameter and returns AuthTokenResponse reference or error.
// For implementation details check LoginExample examples/login.go
func (c *AuthorizerClient) Login(req *LoginInput) (*AuthTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`mutation login($data: LoginInput!) { login(params: $data) { %s } }`, AuthTokenResponseFragment),
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
