package authorizer

import (
	"encoding/json"
	"fmt"
)

// VerifyOTPInput defines attributes for verify_otp request
type VerifyOTPInput struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

// VerifyOTP is method attached to AuthorizerClient.
// It performs verify_otp mutation on authorizer instance.
// It returns AuthTokenResponse reference or error.
// For implementation details check VerifyOTPExample examples/verify_otp.go
func (c *AuthorizerClient) VerifyOTP(req *VerifyOTPInput) (*AuthTokenResponse, error) {
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
