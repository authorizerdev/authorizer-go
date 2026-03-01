package authorizer

import (
	"encoding/json"
)

// ForgotPasswordRequest defines attributes for forgot_password request
type ForgotPasswordRequest struct {
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	State       *string `json:"state,omitempty"`
	RedirectURI *string `json:"redirect_uri,omitempty"`
}

// ForgotPasswordInput is deprecated: Use ForgotPasswordRequest instead
type ForgotPasswordInput = ForgotPasswordRequest

// ForgotPassword is method attached to AuthorizerClient.
// It performs forgot_password mutation on authorizer instance.
// It takes ForgotPasswordRequest reference as parameter and returns ForgotPasswordResponse reference or error.
func (c *AuthorizerClient) ForgotPassword(req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
	if req.State == nil || StringValue(req.State) == "" {
		// generate random state
		req.State = NewStringRef(EncodeB64(CreateRandomString()))
	}

	if req.RedirectURI == nil || StringValue(req.RedirectURI) == "" {
		req.RedirectURI = NewStringRef(c.RedirectURL)
	}

	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `mutation forgotPassword($data: ForgotPasswordRequest!) { forgot_password(params: $data) { message should_show_mobile_otp_screen } }`,
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*ForgotPasswordResponse
	json.Unmarshal(bytesData, &res)

	return res["forgot_password"], nil
}
