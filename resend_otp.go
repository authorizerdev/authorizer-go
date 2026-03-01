package authorizer

import (
	"encoding/json"
)

// ResendOTPRequest defines attributes for resend_otp request
type ResendOTPRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	State       *string `json:"state,omitempty"`
}

// ResendOTPInput is deprecated: Use ResendOTPRequest instead
type ResendOTPInput = ResendOTPRequest

// ResendOTP is method attached to AuthorizerClient.
// It performs resend_otp mutation on authorizer instance.
// It takes ResendOTPRequest reference as parameter and returns Response reference or error.
func (c *AuthorizerClient) ResendOTP(req *ResendOTPRequest) (*Response, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `mutation resendOtp($data: ResendOTPRequest!) { resend_otp(params: $data) { message }}`,
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*Response
	json.Unmarshal(bytesData, &res)

	return res["resend_otp"], nil
}
