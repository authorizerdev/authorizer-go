package authorizer

import (
	"encoding/json"
)

// ResetPasswordInput defines attributes for reset_password request
type ResetPasswordInput struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// ResetPassword is method attached to AuthorizerClient.
// It performs resend_otp mutation on authorizer instance.
// It takes ResetPasswordInput reference as parameter and returns Response reference or error.
// For implementation details check ResetPasswordExample examples/resent_password.go
func (c *AuthorizerClient) ResetPassword(req *ResetPasswordInput) (*Response, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `mutation resetPassword($data: ResetPasswordInput!) {	reset_password(params: $data) { message } }`,
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*Response
	json.Unmarshal(bytesData, &res)

	return res["reset_password"], nil
}
