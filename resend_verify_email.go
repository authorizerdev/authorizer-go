package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// ResendVerifyEmailRequest defines attributes for resend_verify_email request
type ResendVerifyEmailRequest struct {
	Email      string  `json:"email"`
	Identifier *string `json:"identifier,omitempty"`
}

// ResendVerifyEmail is method attached to AuthorizerClient.
// It performs resend_verify_email mutation on authorizer instance.
// It takes ResendVerifyEmailRequest reference as parameter and returns Response reference or error.
func (c *AuthorizerClient) ResendVerifyEmail(req *ResendVerifyEmailRequest) (*Response, error) {
	var res Response
	err := c.execute(methodSpec{
		name: "ResendVerifyEmail",
		graphql: &GraphQLRequest{
			Query:     `mutation resendVerifyEmail($data: ResendVerifyEmailRequest!) { resend_verify_email(params: $data) { message }}`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "resend_verify_email",
		restMethod:   http.MethodPost,
		restPath:     "/v1/resend_verify_email",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.ResendVerifyEmailResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.ResendVerifyEmailRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.ResendVerifyEmail(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
