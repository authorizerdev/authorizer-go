package authorizer

import (
	"encoding/json"
	"fmt"
)

// SessionQueryInput defines attributes for session query request
type SessionQueryInput struct {
	Roles []*string `json:"roles"`
}

// TODO: session currently works with cookie only use this function once the flow is fixed in authorizer core.

// GetSession is method attached to AuthorizerClient.
// It performs session query on authorizer instance.
// It returns User reference or error.
// For implementation details check GetSessionExample examples/get_session.go
func (c *AuthorizerClient) GetSession(req *SessionQueryInput, headers map[string]string) (*AuthTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`query getSession($data: SessionQueryInput){session(params: $data) { %s } }`, AuthTokenResponseFragment),
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
