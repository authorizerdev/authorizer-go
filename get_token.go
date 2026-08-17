package authorizer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GetTokenRequest defines attributes for token request. Only the set (non-nil)
// parameters are sent to the server.
type GetTokenRequest struct {
	Code         *string `json:"code"`
	GrantType    *string `json:"grant_type"`
	RefreshToken *string `json:"refresh_token"`
	// CodeVerifier is the PKCE verifier for public clients using the
	// authorization_code grant (RFC 7636).
	CodeVerifier *string `json:"code_verifier"`
	// ClientSecret authenticates confidential clients (client_credentials grant).
	ClientSecret *string `json:"client_secret"`
	// Scope is the space-delimited OAuth2 scope parameter (RFC 6749 §3.3).
	Scope *string `json:"scope"`
	// ClientAssertion / ClientAssertionType carry the RFC 7523 JWT-bearer client
	// credential (secretless workload-identity path).
	ClientAssertion     *string `json:"client_assertion"`
	ClientAssertionType *string `json:"client_assertion_type"`
	// RFC 8693 token-exchange parameters. SubjectToken carries the authority
	// being exercised; ActorToken carries the acting party.
	SubjectToken     *string `json:"subject_token"`
	SubjectTokenType *string `json:"subject_token_type"`
	ActorToken       *string `json:"actor_token"`
	ActorTokenType   *string `json:"actor_token_type"`
	// Resource is the RFC 8707 resource indicator (target audience).
	Resource *string `json:"resource"`
}

// TokenQueryInput is deprecated: Use GetTokenRequest instead
type TokenQueryInput = GetTokenRequest

// TokenResponse defines attributes for token request
type TokenResponse struct {
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    int64   `json:"expires_in"`
	IdToken      string  `json:"id_token"`
	RefreshToken *string `json:"refresh_token"`
	Scope        string  `json:"scope"`
	// IssuedTokenType is set for RFC 8693 token-exchange responses.
	IssuedTokenType string `json:"issued_token_type"`
}

// oauthError is the RFC 6749 §5.2 error response returned by /oauth/token.
type oauthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GetToken is method attached to AuthorizerClient.
// It performs a form-encoded POST to `/oauth/token` on the authorizer instance
// (OAuth endpoints stay REST regardless of the client's selected protocol).
// It returns TokenResponse reference or error.
func (c *AuthorizerClient) GetToken(req *GetTokenRequest) (*TokenResponse, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	grantType := StringValue(req.GrantType)
	if grantType == "" {
		grantType = GrantTypeAuthorizationCode
	}

	if grantType == GrantTypeRefreshToken && req.RefreshToken == nil {
		return nil, errors.New("invalid refresh token")
	}
	if grantType == GrantTypeAuthorizationCode && req.Code == nil {
		return nil, errors.New("invalid code")
	}

	form := url.Values{}
	form.Set("grant_type", grantType)
	form.Set("client_id", c.ClientID)
	setForm := func(key string, val *string) {
		if val != nil {
			form.Set(key, *val)
		}
	}
	setForm("code", req.Code)
	setForm("code_verifier", req.CodeVerifier)
	setForm("refresh_token", req.RefreshToken)
	setForm("client_secret", req.ClientSecret)
	setForm("scope", req.Scope)
	setForm("client_assertion", req.ClientAssertion)
	setForm("client_assertion_type", req.ClientAssertionType)
	setForm("subject_token", req.SubjectToken)
	setForm("subject_token_type", req.SubjectTokenType)
	setForm("actor_token", req.ActorToken)
	setForm("actor_token_type", req.ActorTokenType)
	setForm("resource", req.Resource)

	httpReq, err := http.NewRequest(http.MethodPost, c.AuthorizerURL+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, val := range c.ExtraHeaders {
		httpReq.Header.Add(key, val)
	}
	// Authorizer's CSRF guard rejects state-changing requests without an Origin
	// or Referer header (see ExecuteGraphQL for the full rationale).
	if httpReq.Header.Get("Origin") == "" {
		if u, err := url.Parse(c.AuthorizerURL); err == nil && u.Scheme != "" && u.Host != "" {
			httpReq.Header.Set("Origin", u.Scheme+"://"+u.Host)
		}
	}

	client := c.HTTPClient()
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		// RFC 6749 §5.2: {"error": "...", "error_description": "..."}.
		var oauthErr oauthError
		if json.Unmarshal(bodyBytes, &oauthErr) == nil && oauthErr.Error != "" {
			if oauthErr.ErrorDescription != "" {
				return nil, fmt.Errorf("%s: %s", oauthErr.Error, oauthErr.ErrorDescription)
			}
			return nil, errors.New(oauthErr.Error)
		}
		return nil, errors.New(http.StatusText(res.StatusCode) + ": " + string(bodyBytes))
	}

	var tokenRes TokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenRes); err != nil {
		return nil, err
	}
	return &tokenRes, nil
}
