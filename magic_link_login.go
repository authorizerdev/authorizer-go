package authorizer

import (
	"encoding/json"
)

// MagicLinkLoginInput defines attributes for magic link login request
type MagicLinkLoginInput struct {
	Email       string    `json:"email"`
	Roles       []*string `json:"roles,omitempty"`
	Scope       []*string `json:"scope,omitempty"`
	State       *string   `json:"state"`
	RedirectURI *string   `json:"redirect_uri"`
}

// MagicLinkLoginInput is method attached to AuthorizerClient.
// It performs magic_link_login mutation on authorizer instance.
// It takes MagicLinkLoginInput reference as parameter and returns AuthTokenResponse reference or error.
// For implementation details check MagicLinkLoginExample examples/magic_link_login.go
func (c *AuthorizerClient) MagicLinkLogin(req *MagicLinkLoginInput) (*Response, error) {
	if req.State == nil || StringValue(req.State) == "" {
		// generate random state
		req.State = NewStringRef(EncodeB64(CreateRandomString()))
	}

	if req.RedirectURI == nil || StringValue(req.RedirectURI) == "" {
		req.RedirectURI = NewStringRef(c.RedirectURL)
	}

	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `mutation magicLinkLogin($data: MagicLinkLoginInput!) { magic_link_login(params: $data) { message }}`,
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*Response
	json.Unmarshal(bytesData, &res)

	return res["magic_link_login"], nil
}
