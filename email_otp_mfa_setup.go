package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-proto-go/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// EmailOtpMfaSetupRequest defines attributes for email_otp_mfa_setup request.
// Email/PhoneNumber are only used in the MFA-session-cookie mode (a caller in
// the withheld first-time-offer state, with no bearer token yet); ignored
// when the caller has a valid bearer token/session, which already identifies
// the user.
type EmailOtpMfaSetupRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

// EmailOtpMfaSetup is method attached to AuthorizerClient.
// It performs email_otp_mfa_setup mutation on authorizer instance, sending a
// one-time code to the caller's own email and creating an unverified
// email-OTP MFA enrollment. Verify the code with VerifyOTP.
// It returns Response reference or error.
func (c *AuthorizerClient) EmailOtpMfaSetup(req *EmailOtpMfaSetupRequest, headers map[string]string) (*Response, error) {
	var res Response
	err := c.execute(methodSpec{
		name: "EmailOtpMfaSetup",
		graphql: &GraphQLRequest{
			Query:     `mutation emailOtpMfaSetup($data: OtpMfaSetupRequest) { email_otp_mfa_setup(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "email_otp_mfa_setup",
		restMethod:   http.MethodPost,
		restPath:     "/v1/email_otp_mfa_setup",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.EmailOtpMfaSetupResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.EmailOtpMfaSetupRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.EmailOtpMfaSetup(ctx, &in)
		},
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
