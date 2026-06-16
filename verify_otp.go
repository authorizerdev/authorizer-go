package authorizer

import (
	"context"
	"fmt"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
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
	var res AuthTokenResponse
	err := c.execute(methodSpec{
		name: "VerifyOTP",
		graphql: &GraphQLRequest{
			Query:     fmt.Sprintf(`mutation verifyOtp($data: VerifyOTPRequest!) { verify_otp(params: $data) { %s }}`, AuthTokenResponseFragment),
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "verify_otp",
		restMethod:   http.MethodPost,
		restPath:     "/v1/verify_otp",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.AuthResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.VerifyOtpRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.VerifyOtp(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
