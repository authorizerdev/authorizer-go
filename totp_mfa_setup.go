package authorizer

import "fmt"

// TotpMfaSetupRequest defines attributes for totp_mfa_setup request. Same
// dual-mode identification semantics as EmailOtpMfaSetupRequest: Email/
// PhoneNumber are only used in the MFA-session-cookie mode; ignored when the
// caller has a valid bearer token/session.
type TotpMfaSetupRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

// TotpMfaSetup is method attached to AuthorizerClient.
// It performs totp_mfa_setup mutation on authorizer instance, generating a
// fresh TOTP secret/QR/recovery-codes for the caller to scan and confirm via
// VerifyOTP(is_totp: true). GraphQL only: the server has no gRPC/REST RPC for
// this operation.
// It returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) TotpMfaSetup(req *TotpMfaSetupRequest, headers map[string]string) (*AuthTokenResponse, error) {
	var res AuthTokenResponse
	err := c.execute(methodSpec{
		name: "TotpMfaSetup",
		graphql: &GraphQLRequest{
			Query:     fmt.Sprintf(`mutation totpMfaSetup($data: OtpMfaSetupRequest) { totp_mfa_setup(params: $data) { %s }}`, AuthTokenResponseFragment),
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "totp_mfa_setup",
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
