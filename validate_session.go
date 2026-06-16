package authorizer

import (
	"context"
	"fmt"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// ValidateSessionRequest defines attributes for validate_session request
type ValidateSessionRequest struct {
	Cookie string    `json:"cookie,omitempty"`
	Roles  []*string `json:"roles,omitempty"`
}

// ValidateSessionInput is deprecated: Use ValidateSessionRequest instead
type ValidateSessionInput = ValidateSessionRequest

// ValidateSessionResponse defines attributes for validate_session response
type ValidateSessionResponse struct {
	IsValid bool  `json:"is_valid"`
	User    *User `json:"user,omitempty"`
}

// ValidateSession is method attached to AuthorizerClient.
// It performs validate_session query on authorizer instance.
// It returns ValidateSessionResponse reference or error.
func (c *AuthorizerClient) ValidateSession(req *ValidateSessionRequest) (*ValidateSessionResponse, error) {
	var res ValidateSessionResponse
	err := c.execute(methodSpec{
		name: "ValidateSession",
		graphql: &GraphQLRequest{
			Query:     fmt.Sprintf(`query validateSession($data: ValidateSessionRequest!){validate_session(params: $data) { is_valid user { %s } } }`, UserFragment),
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "validate_session",
		restMethod:   http.MethodPost,
		restPath:     "/v1/validate_session",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.ValidateSessionResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.ValidateSessionRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.ValidateSession(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
