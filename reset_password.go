package authorizer

import (
	"encoding/json"
)

// ResetPasswordRequest defines attributes for reset_password request
type ResetPasswordRequest struct {
	Token           *string `json:"token,omitempty"`
	Password        string  `json:"password"`
	ConfirmPassword string  `json:"confirm_password"`
	OTP             *string `json:"otp,omitempty"`
	PhoneNumber     *string `json:"phone_number,omitempty"`
}

// ResetPasswordInput is deprecated: Use ResetPasswordRequest instead
type ResetPasswordInput = ResetPasswordRequest

// ResetPassword is method attached to AuthorizerClient.
// It performs reset_password mutation on authorizer instance.
// It takes ResetPasswordRequest reference as parameter and returns Response reference or error.
func (c *AuthorizerClient) ResetPassword(req *ResetPasswordRequest) (*Response, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `mutation resetPassword($data: ResetPasswordRequest!) {	reset_password(params: $data) { message } }`,
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
