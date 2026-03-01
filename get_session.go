package authorizer

import (
	"encoding/json"
	"fmt"
)

// SessionQueryRequest defines attributes for session query request
type SessionQueryRequest struct {
	Roles []*string `json:"roles"`
	Scope []*string `json:"scope,omitempty"`
}

// SessionQueryInput is deprecated: Use SessionQueryRequest instead
type SessionQueryInput = SessionQueryRequest

// GetSession is method attached to AuthorizerClient.
// It performs session query on authorizer instance.
// It returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) GetSession(req *SessionQueryRequest, headers map[string]string) (*AuthTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`query getSession($data: SessionQueryRequest){session(params: $data) { %s } }`, AuthTokenResponseFragment),
		Variables: map[string]interface{}{
			"data": req,
		},
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]*AuthTokenResponse
	json.Unmarshal(bytesData, &res)

	return res["session"], nil
}
