package authorizer

import (
	"context"
	"fmt"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// GetProfile is method attached to AuthorizerClient.
// It performs profile query on authorizer instance.
// It returns User reference or error.
// For implementation details check GetProfileExample examples/get_profile.go
func (c *AuthorizerClient) GetProfile(headers map[string]string) (*User, error) {
	var res User
	err := c.execute(methodSpec{
		name: "GetProfile",
		graphql: &GraphQLRequest{
			Query:     fmt.Sprintf(`query {	profile { %s } }`, UserFragment),
			Variables: nil,
		},
		graphqlField: "profile",
		restMethod:   http.MethodGet,
		restPath:     "/v1/profile",
		restBody:     nil,
		restResp:     func() proto.Message { return &authorizerv1.User{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			return cli.Profile(ctx, &authorizerv1.ProfileRequest{})
		},
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
