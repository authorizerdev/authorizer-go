package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
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
	var res Response
	err := c.execute(methodSpec{
		name: "ResendOTP",
		graphql: &GraphQLRequest{
			Query:     `mutation resendOtp($data: ResendOTPRequest!) { resend_otp(params: $data) { message }}`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "resend_otp",
		restMethod:   http.MethodPost,
		restPath:     "/v1/resend_otp",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.ResendOtpResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.ResendOtpRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.ResendOtp(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
