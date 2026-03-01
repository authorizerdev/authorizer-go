package authorizer

import (
	"encoding/json"
	"fmt"
)

// ValidateSessionRequest defines attributes for validate_session request
type ValidateSessionRequest struct {
	Cookie string    `json:"cookie,omitempty"`
	Roles  []*string `json:"roles,omitempty"`
}

// ValidateSessionInput is deprecated: Use ValidateSessionRequest instead
type ValidateSessionInput = ValidateSessionRequest

// ValidateSessionResponse defines attributes for validate_session response
type ValidateSessionResponse struct {
	IsValid bool  `json:"is_valid"`
	User    *User `json:"user,omitempty"`
}

// ValidateSession is method attached to AuthorizerClient.
// It performs validate_session query on authorizer instance.
// It returns ValidateSessionResponse reference or error.
func (c *AuthorizerClient) ValidateSession(req *ValidateSessionRequest) (*ValidateSessionResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`query validateSession($data: ValidateSessionRequest!){validate_session(params: $data) { is_valid user { %s } } }`, UserFragment),
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*ValidateSessionResponse
	json.Unmarshal(bytesData, &res)

	return res["validate_session"], nil
}
