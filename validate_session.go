package authorizer

import "encoding/json"

// ValidateSessionInput defines attributes for validate_session request
type ValidateSessionInput struct {
	Cookie string    `json:"cookie,omitempty"`
	Roles  []*string `json:"roles,omitempty"`
}

// ValidateSessionResponse defines attributes for validate_session response

type ValidateSessionResponse struct {
	IsValid bool  `json:"is_valid"`
	User    *User `json:"user,omitempty"`
}

// ValidateSession is method attached to AuthorizerClient.
// It performs validate_session query on authorizer instance.
// It returns ValidateSessionResponse reference or error.
// For implementation details check ValidateSessionExample examples/validate_session.go
func (c *AuthorizerClient) ValidateSession(req *ValidateSessionInput) (*ValidateSessionResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `query validateSession($data: ValidateSessionInput!){validate_session(params: $data) { is_valid user } }`,
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
