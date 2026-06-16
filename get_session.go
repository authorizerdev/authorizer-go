package authorizer

import (
	"context"
	"fmt"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// SessionQueryRequest defines attributes for session query request
type SessionQueryRequest struct {
	Roles []*string `json:"roles"`
	Scope []*string `json:"scope,omitempty"`
}

// SessionQueryInput is deprecated: Use SessionQueryRequest instead
type SessionQueryInput = SessionQueryRequest

// GetSession is method attached to AuthorizerClient.
// It performs session query on authorizer instance.
// It returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) GetSession(req *SessionQueryRequest, headers map[string]string) (*AuthTokenResponse, error) {
	var res AuthTokenResponse
	err := c.execute(methodSpec{
		name: "GetSession",
		graphql: &GraphQLRequest{
			Query:     fmt.Sprintf(`query getSession($data: SessionQueryRequest){session(params: $data) { %s } }`, AuthTokenResponseFragment),
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "session",
		restMethod:   http.MethodPost,
		restPath:     "/v1/session",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.AuthResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.SessionRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.Session(ctx, &in)
		},
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
