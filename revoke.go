package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// RevokeRequest defines attributes for the revoke mutation/RPC (distinct from
// the OAuth2-standard RevokeToken, which targets /oauth/revoke directly).
type RevokeRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Revoke is method attached to AuthorizerClient.
// It performs the revoke mutation/RPC (/v1/revoke) on authorizer instance
// over the client's selected protocol (graphql, rest or grpc). For the
// OAuth2-standard token revocation endpoint (/oauth/revoke), use RevokeToken.
// It returns Response reference or error.
func (c *AuthorizerClient) Revoke(req *RevokeRequest) (*Response, error) {
	var res Response
	err := c.execute(methodSpec{
		name: "Revoke",
		graphql: &GraphQLRequest{
			Query:     `mutation revoke($data: OAuthRevokeRequest!) { revoke(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "revoke",
		restMethod:   http.MethodPost,
		restPath:     "/v1/revoke",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.RevokeResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.RevokeRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.Revoke(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
