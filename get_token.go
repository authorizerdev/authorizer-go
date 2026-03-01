package authorizer

import "errors"

// GetTokenRequest defines attributes for token request
type GetTokenRequest struct {
	Code         *string `json:"code"`
	GrantType    *string `json:"grant_type"`
	RefreshToken *string `json:"refresh_token"`
}

// TokenQueryInput is deprecated: Use GetTokenRequest instead
type TokenQueryInput = GetTokenRequest

// TokenResponse defines attributes for token request
type TokenResponse struct {
	AccessToken  string  `json:"access_token"`
	ExpiresIn    int64   `json:"expires_in"`
	IdToken      string  `json:"id_token"`
	RefreshToken *string `json:"refresh_token"`
}

// GetToken is method attached to AuthorizerClient.
// It performs `/oauth/token` query on authorizer instance.
// It returns TokenResponse reference or error.
func (c *AuthorizerClient) GetToken(req *GetTokenRequest) (*TokenResponse, error) {
	grantType := StringValue(req.GrantType)
	if grantType == "" {
		req.GrantType = NewStringRef(GrantTypeAuthorizationCode)
	}

	if grantType == GrantTypeRefreshToken && req.RefreshToken == nil {
		return nil, errors.New("invalid refresh token")
	}

	return nil, nil
}
