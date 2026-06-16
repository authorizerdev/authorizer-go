package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
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
	var res Response
	err := c.execute(methodSpec{
		name: "ResetPassword",
		graphql: &GraphQLRequest{
			Query:     `mutation resetPassword($data: ResetPasswordRequest!) {	reset_password(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "reset_password",
		restMethod:   http.MethodPost,
		restPath:     "/v1/reset_password",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.ResetPasswordResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.ResetPasswordRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.ResetPassword(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
