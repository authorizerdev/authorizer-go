package authorizer

import (
	"encoding/json"
	"fmt"
)

// VerifyOTPRequest defines attributes for verify_otp request
type VerifyOTPRequest struct {
	Email       *string `json:"email"`
	OTP         string  `json:"otp"`
	PhoneNumber *string `json:"phone_number"`
	IsTotp      *bool   `json:"is_totp,omitempty"`
	State       *string `json:"state,omitempty"`
}

// VerifyOTPInput is deprecated: Use VerifyOTPRequest instead
type VerifyOTPInput = VerifyOTPRequest

// VerifyOTP is method attached to AuthorizerClient.
// It performs verify_otp mutation on authorizer instance.
// It returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) VerifyOTP(req *VerifyOTPRequest) (*AuthTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`mutation verifyOtp($data: VerifyOTPRequest!) { verify_otp(params: $data) { %s }}`, AuthTokenResponseFragment),
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*AuthTokenResponse
	json.Unmarshal(bytesData, &res)

	return res["verify_otp"], nil
}
