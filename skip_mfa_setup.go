package authorizer

import (
	"context"
	"fmt"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// SkipMfaSetupRequest defines attributes for skip_mfa_setup request
type SkipMfaSetupRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	State       *string `json:"state,omitempty"`
}

// SkipMfaSetup is method attached to AuthorizerClient.
// It performs skip_mfa_setup mutation on authorizer instance, completing an
// in-progress, token-withheld MFA offer by recording that the caller
// explicitly declined it, then issuing the access token that was withheld.
// Fails if MFA is organization-enforced (enforcement is never skippable).
// It returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) SkipMfaSetup(req *SkipMfaSetupRequest) (*AuthTokenResponse, error) {
	var res AuthTokenResponse
	err := c.execute(methodSpec{
		name: "SkipMfaSetup",
		graphql: &GraphQLRequest{
			Query:     fmt.Sprintf(`mutation skipMfaSetup($data: SkipMfaSetupRequest!) { skip_mfa_setup(params: $data) { %s }}`, AuthTokenResponseFragment),
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "skip_mfa_setup",
		restMethod:   http.MethodPost,
		restPath:     "/v1/skip_mfa_setup",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.AuthResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.SkipMfaSetupRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.SkipMfaSetup(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
