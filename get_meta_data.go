package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// MetaDataResponse defines attributes for MetaData response query
type MetaDataResponse struct {
	Version                            string `json:"version"`
	ClientID                           string `json:"client_id"`
	IsGoogleLoginEnabled               bool   `json:"is_google_login_enabled"`
	IsFacebookLoginEnabled             bool   `json:"is_facebook_login_enabled"`
	IsGithubLoginEnabled               bool   `json:"is_github_login_enabled"`
	IsLinkedinLoginEnabled             bool   `json:"is_linkedin_login_enabled"`
	IsAppleLoginEnabled                bool   `json:"is_apple_login_enabled"`
	IsTwitterLoginEnabled              bool   `json:"is_twitter_login_enabled"`
	IsDiscordLoginEnabled              bool   `json:"is_discord_login_enabled"`
	IsMicrosoftLoginEnabled            bool   `json:"is_microsoft_login_enabled"`
	IsTwitchLoginEnabled               bool   `json:"is_twitch_login_enabled"`
	IsRobloxLoginEnabled               bool   `json:"is_roblox_login_enabled"`
	IsEmailVerificationEnabled         bool   `json:"is_email_verification_enabled"`
	IsBasicAuthenticationEnabled       bool   `json:"is_basic_authentication_enabled"`
	IsMagicLinkLoginEnabled            bool   `json:"is_magic_link_login_enabled"`
	IsSignUpEnabled                    bool   `json:"is_sign_up_enabled"`
	IsStrongPasswordEnabled            bool   `json:"is_strong_password_enabled"`
	IsMultiFactorAuthEnabled           bool   `json:"is_multi_factor_auth_enabled"`
	IsMobileBasicAuthenticationEnabled bool   `json:"is_mobile_basic_authentication_enabled"`
	IsPhoneVerificationEnabled         bool   `json:"is_phone_verification_enabled"`
}

// GetMetaData is method attached to AuthorizerClient.
// It performs meta query on authorizer instance.
// It returns MetaResponse reference or error.
func (c *AuthorizerClient) GetMetaData() (*MetaDataResponse, error) {
	var res MetaDataResponse
	err := c.execute(methodSpec{
		name: "GetMetaData",
		graphql: &GraphQLRequest{
			Query:     `query { meta { version client_id is_google_login_enabled is_facebook_login_enabled is_github_login_enabled is_linkedin_login_enabled is_apple_login_enabled is_twitter_login_enabled is_discord_login_enabled is_microsoft_login_enabled is_twitch_login_enabled is_roblox_login_enabled is_email_verification_enabled is_basic_authentication_enabled is_magic_link_login_enabled is_sign_up_enabled is_strong_password_enabled is_multi_factor_auth_enabled is_mobile_basic_authentication_enabled is_phone_verification_enabled } }`,
			Variables: nil,
		},
		graphqlField: "meta",
		restMethod:   http.MethodGet,
		restPath:     "/v1/meta",
		restBody:     nil,
		restResp:     func() proto.Message { return &authorizerv1.Meta{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			return cli.Meta(ctx, &authorizerv1.MetaRequest{})
		},
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
