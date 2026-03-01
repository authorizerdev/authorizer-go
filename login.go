package authorizer

import (
	"encoding/json"
	"fmt"
)

// LoginRequest defines attributes for login request
type LoginRequest struct {
	Email       *string   `json:"email,omitempty"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Password    string    `json:"password"`
	Roles       []*string `json:"roles,omitempty"`
	Scope       []*string `json:"scope,omitempty"`
	State       *string   `json:"state,omitempty"`
}

// LoginInput is deprecated: Use LoginRequest instead
type LoginInput = LoginRequest

// Login is method attached to AuthorizerClient.
// It performs login mutation on authorizer instance.
// It takes LoginRequest reference as parameter and returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) Login(req *LoginRequest) (*AuthTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`mutation login($data: LoginRequest!) { login(params: $data) { %s } }`, AuthTokenResponseFragment),
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
