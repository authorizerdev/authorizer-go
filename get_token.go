package authorizer

import "errors"

// TokenQueryInput defines attributes for token request
type TokenQueryInput struct {
	Code         *string `json:"code"`
	GrantType    *string `json:"grant_type"`
	RefreshToken *string `json:"refresh_token"`
}

// TODO check if we can use oauth get token from backend

// TokenResponse defines attributes for token request
type TokenResponse struct {
	AccessToken  string  `json:"access_token"`
	ExpiresIn    int64   `json:"expires_in"`
	IdToken      string  `json:"id_token"`
	RefreshToken *string `json:"refresh_token"`
}

// GetToken is method attached to AuthorizerClient.
// It performs `/oauth/token` query on authorizer instance.
// It returns User reference or error.
// For implementation details check GetTokenExample examples/get_token.go
func (c *AuthorizerClient) GetToken(req *TokenQueryInput) (*TokenResponse, error) {
	grantType := StringValue(req.GrantType)
	if grantType == "" {
		req.GrantType = NewStringRef(GrantTypeAuthorizationCode)
	}

	if grantType == GrantTypeRefreshToken && req.RefreshToken == nil {
		return nil, errors.New("invalid refresh token")
	}

	return nil, nil
}
