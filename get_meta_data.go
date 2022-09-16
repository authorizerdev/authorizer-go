package authorizer

import (
	"encoding/json"
)

// MetaDataResponse defines attributes for MetaData response query
type MetaDataResponse struct {
	Version                      string `json:"version"`
	ClientID                     string `json:"client_id"`
	IsGoogleLoginEnabled         bool   `json:"is_google_login_enabled"`
	IsFacebookLoginEnabled       bool   `json:"is_facebook_login_enabled"`
	IsGithubLoginEnabled         bool   `json:"is_github_login_enabled"`
	IsLinkedinLoginEnabled       bool   `json:"is_linkedin_login_enabled"`
	IsAppleLoginEnabled          bool   `json:"is_apple_login_enabled"`
	IsTwitterLoginEnabled        bool   `json:"is_twitter_login_enabled"`
	IsEmailVerificationEnabled   bool   `json:"is_email_verification_enabled"`
	IsBasicAuthenticationEnabled bool   `json:"is_basic_authentication_enabled"`
	IsMagicLinkLoginEnabled      bool   `json:"is_magic_link_login_enabled"`
	IsSignUpEnabled              bool   `json:"is_sign_up_enabled"`
	IsStrongPasswordEnabled      bool   `json:"is_strong_password_enabled"`
}

// GetMetaData is method attached to AuthorizerClient.
// It performs meta query on authorizer instance.
// It returns MetaResponse reference or error.
// For implementation details check GetMetadataExample examples/get_meta_data.go
func (c *AuthorizerClient) GetMetaData() (*MetaDataResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query:     `query { meta { version client_id is_google_login_enabled is_facebook_login_enabled is_github_login_enabled is_linkedin_login_enabled is_apple_login_enabled is_twitter_login_enabled is_email_verification_enabled is_basic_authentication_enabled is_magic_link_login_enabled is_sign_up_enabled is_strong_password_enabled } }`,
		Variables: nil,
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*MetaDataResponse
	json.Unmarshal(bytesData, &res)

	return res["meta"], nil
}
