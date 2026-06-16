package authorizer

import (
	"context"
	"fmt"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// LoginRequest defines attributes for login request
type LoginRequest struct {
	Email       *string   `json:"email,omitempty"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Password    string    `json:"password"`
	Roles       []*string `json:"roles,omitempty"`
	Scope       []*string `json:"scope,omitempty"`
	State       *string   `json:"state,omitempty"`
}

// LoginInput is deprecated: Use LoginRequest instead
type LoginInput = LoginRequest

// Login is method attached to AuthorizerClient.
// It performs login on the authorizer instance over the client's selected
// protocol (graphql, rest or grpc).
// It takes LoginRequest reference as parameter and returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) Login(req *LoginRequest) (*AuthTokenResponse, error) {
	var res AuthTokenResponse
	err := c.execute(methodSpec{
		name: "Login",
		graphql: &GraphQLRequest{
			Query:     fmt.Sprintf(`mutation login($data: LoginRequest!) { login(params: $data) { %s } }`, AuthTokenResponseFragment),
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "login",
		restMethod:   http.MethodPost,
		restPath:     "/v1/login",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.AuthResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.LoginRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.Login(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
