package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// TODO does not work without cookie

// Logout is method attached to AuthorizerClient.
// It performs Logout mutation on authorizer instance.
// It takes LogoutInput reference as parameter and returns Response reference or error.
// For implementation details check LogoutExample examples/Logout.go
func (c *AuthorizerClient) Logout(headers map[string]string) (*Response, error) {
	var res Response
	err := c.execute(methodSpec{
		name: "Logout",
		graphql: &GraphQLRequest{
			Query:     `mutation { logout { message } }`,
			Variables: nil,
		},
		graphqlField: "logout",
		restMethod:   http.MethodPost,
		restPath:     "/v1/logout",
		restBody:     nil,
		restResp:     func() proto.Message { return &authorizerv1.LogoutResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			return cli.Logout(ctx, &authorizerv1.LogoutRequest{})
		},
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
