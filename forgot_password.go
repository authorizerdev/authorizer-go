package authorizer

import (
	"encoding/json"
)

// ForgotPasswordInput defines attributes for forgot_password request
type ForgotPasswordInput struct {
	Email       string  `json:"email"`
	State       *string `json:"state,omitempty"`
	RedirectURI *string `json:"redirect_uri,omitempty"`
}

// ForgotPassword is method attached to AuthorizerClient.
// It performs forgot_password mutation on authorizer instance.
// It takes ForgotPasswordInput reference as parameter and returns Response reference or error.
// For implementation details check ForgotPasswordInputExample in examples/forgot_password.go
func (c *AuthorizerClient) ForgotPassword(req *ForgotPasswordInput) (*Response, error) {
	if req.State == nil || StringValue(req.State) == "" {
		// generate random state
		req.State = NewStringRef(EncodeB64(CreateRandomString()))
	}

	if req.RedirectURI == nil || StringValue(req.RedirectURI) == "" {
		req.RedirectURI = NewStringRef(c.RedirectURL)
	}

	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `mutation forgotPassword($data: ForgotPasswordInput!) { forgot_password(params: $data) { message } }`,
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*Response
	json.Unmarshal(bytesData, &res)

	return res["forgot_password"], nil
}
