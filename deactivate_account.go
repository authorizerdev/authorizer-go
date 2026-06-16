package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// DeactivateAccount is method attached to AuthorizerClient.
// It performs deactivate_account mutation on the authorizer instance over the
// client's selected protocol (graphql, rest or grpc).
// It returns Response reference or error.
// For implementation details check DeactivateAccountExample examples/deactivate_account.go
func (c *AuthorizerClient) DeactivateAccount(headers map[string]string) (*Response, error) {
	var res Response
	err := c.execute(methodSpec{
		name: "DeactivateAccount",
		graphql: &GraphQLRequest{
			Query:     `mutation deactivateAccount { deactivate_account { message } }`,
			Variables: nil,
		},
		graphqlField: "deactivate_account",
		restMethod:   http.MethodPost,
		restPath:     "/v1/deactivate_account",
		restBody:     nil,
		restResp:     func() proto.Message { return &authorizerv1.DeactivateAccountResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			return cli.DeactivateAccount(ctx, &authorizerv1.DeactivateAccountRequest{})
		},
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
