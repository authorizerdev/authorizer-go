package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// ForgotPasswordRequest defines attributes for forgot_password request
type ForgotPasswordRequest struct {
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	State       *string `json:"state,omitempty"`
	RedirectURI *string `json:"redirect_uri,omitempty"`
}

// ForgotPasswordInput is deprecated: Use ForgotPasswordRequest instead
type ForgotPasswordInput = ForgotPasswordRequest

// ForgotPassword is method attached to AuthorizerClient.
// It performs forgot_password mutation on authorizer instance.
// It takes ForgotPasswordRequest reference as parameter and returns ForgotPasswordResponse reference or error.
func (c *AuthorizerClient) ForgotPassword(req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
	if req.State == nil || StringValue(req.State) == "" {
		// generate random state
		req.State = NewStringRef(EncodeB64(CreateRandomString()))
	}

	if req.RedirectURI == nil || StringValue(req.RedirectURI) == "" {
		req.RedirectURI = NewStringRef(c.RedirectURL)
	}

	var res ForgotPasswordResponse
	err := c.execute(methodSpec{
		name: "ForgotPassword",
		graphql: &GraphQLRequest{
			Query:     `mutation forgotPassword($data: ForgotPasswordRequest!) { forgot_password(params: $data) { message should_show_mobile_otp_screen } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "forgot_password",
		restMethod:   http.MethodPost,
		restPath:     "/v1/forgot_password",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.ForgotPasswordResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.ForgotPasswordRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.ForgotPassword(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
