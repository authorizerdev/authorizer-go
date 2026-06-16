package authorizer

import (
	"context"
	"fmt"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// VerifyEmailRequest defines attributes for verify_email request
type VerifyEmailRequest struct {
	Token string  `json:"token"`
	State *string `json:"state,omitempty"`
}

// VerifyEmailInput is deprecated: Use VerifyEmailRequest instead
type VerifyEmailInput = VerifyEmailRequest

// VerifyEmail is method attached to AuthorizerClient.
// It performs verify_email mutation on authorizer instance.
// It returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) VerifyEmail(req *VerifyEmailRequest) (*AuthTokenResponse, error) {
	var res AuthTokenResponse
	err := c.execute(methodSpec{
		name: "VerifyEmail",
		graphql: &GraphQLRequest{
			Query:     fmt.Sprintf(`mutation verifyEmail($data: VerifyEmailRequest!) { verify_email(params: $data) { %s}}`, AuthTokenResponseFragment),
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "verify_email",
		restMethod:   http.MethodPost,
		restPath:     "/v1/verify_email",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.AuthResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.VerifyEmailRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.VerifyEmail(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
