package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// SmsOtpMfaSetupRequest defines attributes for sms_otp_mfa_setup request. Same
// dual-mode identification semantics as EmailOtpMfaSetupRequest.
type SmsOtpMfaSetupRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

// SmsOtpMfaSetup is method attached to AuthorizerClient.
// It performs sms_otp_mfa_setup mutation on authorizer instance, sending a
// one-time code to the caller's own phone number and creating an unverified
// SMS-OTP MFA enrollment. Verify the code with VerifyOTP.
// It returns Response reference or error.
func (c *AuthorizerClient) SmsOtpMfaSetup(req *SmsOtpMfaSetupRequest, headers map[string]string) (*Response, error) {
	var res Response
	err := c.execute(methodSpec{
		name: "SmsOtpMfaSetup",
		graphql: &GraphQLRequest{
			Query:     `mutation smsOtpMfaSetup($data: OtpMfaSetupRequest) { sms_otp_mfa_setup(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "sms_otp_mfa_setup",
		restMethod:   http.MethodPost,
		restPath:     "/v1/sms_otp_mfa_setup",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.SmsOtpMfaSetupResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.SmsOtpMfaSetupRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.SmsOtpMfaSetup(ctx, &in)
		},
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
