package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// UpdateProfileRequest defines attributes for update_profile request
type UpdateProfileRequest struct {
	Email                    *string                `json:"email,omitempty"`
	NewPassword              *string                `json:"new_password,omitempty"`
	ConfirmNewPassword       *string                `json:"confirm_new_password,omitempty"`
	OldPassword              *string                `json:"old_password,omitempty"`
	GivenName                *string                `json:"given_name,omitempty"`
	FamilyName               *string                `json:"family_name,omitempty"`
	MiddleName               *string                `json:"middle_name,omitempty"`
	NickName                 *string                `json:"nick_name,omitempty"`
	Picture                  *string                `json:"picture,omitempty"`
	Gender                   *string                `json:"gender,omitempty"`
	BirthDate                *string                `json:"birthdate,omitempty"`
	PhoneNumber              *string                `json:"phone_number,omitempty"`
	Roles                    []*string              `json:"roles,omitempty"`
	Scope                    []*string              `json:"scope,omitempty"`
	RedirectURI              *string                `json:"redirect_uri,omitempty"`
	IsMultiFactorAuthEnabled *bool                  `json:"is_multi_factor_auth_enabled,omitempty"`
	AppData                  map[string]interface{} `json:"app_data,omitempty"`
}

// UpdateProfileInput is deprecated: Use UpdateProfileRequest instead
type UpdateProfileInput = UpdateProfileRequest

// UpdateProfile is method attached to AuthorizerClient.
// It performs update_profile mutation on authorizer instance.
// It returns Response reference or error.
func (c *AuthorizerClient) UpdateProfile(req *UpdateProfileRequest, headers map[string]string) (*Response, error) {
	var res Response
	err := c.execute(methodSpec{
		name: "UpdateProfile",
		graphql: &GraphQLRequest{
			Query:     `mutation updateProfile($data: UpdateProfileRequest!) {	update_profile(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "update_profile",
		restMethod:   http.MethodPost,
		restPath:     "/v1/update_profile",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.UpdateProfileResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.UpdateProfileRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.UpdateProfile(ctx, &in)
		},
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
