package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-proto-go/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// LockMfaRequest defines attributes for lock_mfa request
type LockMfaRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

// LockMfa is method attached to AuthorizerClient.
// It performs lock_mfa mutation on authorizer instance, recording that the
// caller lost access to their only MFA factor(s). Only allowed when the
// caller has no verified email/SMS OTP fallback enrolled. Does not issue a
// token; the account requires admin recovery afterward.
// It returns Response reference or error.
func (c *AuthorizerClient) LockMfa(req *LockMfaRequest) (*Response, error) {
	var res Response
	err := c.execute(methodSpec{
		name: "LockMfa",
		graphql: &GraphQLRequest{
			Query:     `mutation lockMfa($data: LockMfaRequest!) { lock_mfa(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "lock_mfa",
		restMethod:   http.MethodPost,
		restPath:     "/v1/lock_mfa",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.LockMfaResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.LockMfaRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.LockMfa(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
