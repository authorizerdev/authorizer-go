package authorizer

import (
	"encoding/json"
)

// ResendOTPInput defines attributes for resend_otp request
type ResendOTPInput struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

// ResendOTP is method attached to AuthorizerClient.
// It performs resend_otp mutation on authorizer instance.
// It takes ResendOTPInput reference as parameter and returns Response reference or error.
// For implementation details check ResendOTPExample examples/resend_otp.go
func (c *AuthorizerClient) ResendOTP(req *ResendOTPInput) (*Response, error) {
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
