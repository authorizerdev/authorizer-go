package authorizer

import "encoding/json"

type TokenType string

const (
	// Token type access_token
	TokenTypeAccessToken TokenType = "access_token"
	// Token type id_token
	TokenTypeIDToken TokenType = "id_token"
	// Token type refresh_token
	TokenTypeRefreshToken TokenType = "refresh_token"
)

// ValidateJWTTokenInput defines attributes for validate_jwt_token request
type ValidateJWTTokenInput struct {
	TokenType TokenType `json:"token_type"`
	Token     string    `json:"token"`
	Roles     []*string `json:"roles,omitempty"`
}

// ValidateJWTTokenResponse defines attributes for validate_jwt_token response

type ValidateJWTTokenResponse struct {
	IsValid bool `json:"is_valid"`
}

// ValidateJWTToken is method attached to AuthorizerClient.
// It performs validate_jwt_token query on authorizer instance.
// It returns ValidateJWTTokenResponse reference or error.
// For implementation details check ValidateJWTTokenExample examples/validate_jwt_token.go
func (c *AuthorizerClient) ValidateJWTToken(req *ValidateJWTTokenInput) (*ValidateJWTTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `query validateJWTToken($data: ValidateJWTTokenInput!){validate_jwt_token(params: $data) { is_valid } }`,
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*ValidateJWTTokenResponse
	json.Unmarshal(bytesData, &res)

	return res["validate_jwt_token"], nil
}
