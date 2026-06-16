package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// MagicLinkLoginRequest defines attributes for magic link login request
type MagicLinkLoginRequest struct {
	Email       string    `json:"email"`
	Roles       []*string `json:"roles,omitempty"`
	Scope       []*string `json:"scope,omitempty"`
	State       *string   `json:"state"`
	RedirectURI *string   `json:"redirect_uri"`
}

// MagicLinkLoginInput is deprecated: Use MagicLinkLoginRequest instead
type MagicLinkLoginInput = MagicLinkLoginRequest

// MagicLinkLogin is method attached to AuthorizerClient.
// It performs magic_link_login mutation on the authorizer instance over the
// client's selected protocol (graphql, rest or grpc).
// It takes MagicLinkLoginRequest reference as parameter and returns Response reference or error.
func (c *AuthorizerClient) MagicLinkLogin(req *MagicLinkLoginRequest) (*Response, error) {
	if req.State == nil || StringValue(req.State) == "" {
		// generate random state
		req.State = NewStringRef(EncodeB64(CreateRandomString()))
	}

	if req.RedirectURI == nil || StringValue(req.RedirectURI) == "" {
		req.RedirectURI = NewStringRef(c.RedirectURL)
	}

	var res Response
	err := c.execute(methodSpec{
		name: "MagicLinkLogin",
		graphql: &GraphQLRequest{
			Query:     `mutation magicLinkLogin($data: MagicLinkLoginRequest!) { magic_link_login(params: $data) { message }}`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "magic_link_login",
		restMethod:   http.MethodPost,
		restPath:     "/v1/magic_link_login",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.MagicLinkLoginResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.MagicLinkLoginRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.MagicLinkLogin(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
