package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

type TokenType string

const (
	// Token type access_token
	TokenTypeAccessToken TokenType = "access_token"
	// Token type id_token
	TokenTypeIDToken TokenType = "id_token"
	// Token type refresh_token
	TokenTypeRefreshToken TokenType = "refresh_token"
)

// ValidateJWTTokenRequest defines attributes for validate_jwt_token request
type ValidateJWTTokenRequest struct {
	TokenType TokenType `json:"token_type"`
	Token     string    `json:"token"`
	Roles     []*string `json:"roles,omitempty"`
}

// ValidateJWTTokenInput is deprecated: Use ValidateJWTTokenRequest instead
type ValidateJWTTokenInput = ValidateJWTTokenRequest

// ValidateJWTTokenResponse defines attributes for validate_jwt_token response
type ValidateJWTTokenResponse struct {
	IsValid bool                   `json:"is_valid"`
	Claims  map[string]interface{} `json:"claims"`
}

// ValidateJWTToken is method attached to AuthorizerClient.
// It performs validate_jwt_token query on authorizer instance.
// It returns ValidateJWTTokenResponse reference or error.
func (c *AuthorizerClient) ValidateJWTToken(req *ValidateJWTTokenRequest) (*ValidateJWTTokenResponse, error) {
	var res ValidateJWTTokenResponse
	err := c.execute(methodSpec{
		name: "ValidateJWTToken",
		graphql: &GraphQLRequest{
			Query:     `query validateJWTToken($data: ValidateJWTTokenRequest!){validate_jwt_token(params: $data) { is_valid claims } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "validate_jwt_token",
		restMethod:   http.MethodPost,
		restPath:     "/v1/validate_jwt_token",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.ValidateJwtTokenResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.ValidateJwtTokenRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.ValidateJwtToken(ctx, &in)
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
