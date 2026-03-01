package authorizer

import (
	"encoding/json"
)

// MagicLinkLoginRequest defines attributes for magic link login request
type MagicLinkLoginRequest struct {
	Email       string    `json:"email"`
	Roles       []*string `json:"roles,omitempty"`
	Scope       []*string `json:"scope,omitempty"`
	State       *string   `json:"state"`
	RedirectURI *string   `json:"redirect_uri"`
}

// MagicLinkLoginInput is deprecated: Use MagicLinkLoginRequest instead
type MagicLinkLoginInput = MagicLinkLoginRequest

// MagicLinkLogin is method attached to AuthorizerClient.
// It performs magic_link_login mutation on authorizer instance.
// It takes MagicLinkLoginRequest reference as parameter and returns Response reference or error.
func (c *AuthorizerClient) MagicLinkLogin(req *MagicLinkLoginRequest) (*Response, error) {
	if req.State == nil || StringValue(req.State) == "" {
		// generate random state
		req.State = NewStringRef(EncodeB64(CreateRandomString()))
	}

	if req.RedirectURI == nil || StringValue(req.RedirectURI) == "" {
		req.RedirectURI = NewStringRef(c.RedirectURL)
	}

	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `mutation magicLinkLogin($data: MagicLinkLoginRequest!) { magic_link_login(params: $data) { message }}`,
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
