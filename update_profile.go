package authorizer

import (
	"encoding/json"
)

// UpdateProfileInput defines attributes for signup request
type UpdateProfileInput struct {
	Email                    *string   `json:"email,omitempty"`
	NewPassword              *string   `json:"new_password,omitempty"`
	ConfirmNewPassword       *string   `json:"confirm_new_password,omitempty"`
	OldPassword              *string   `json:"old_password,omitempty"`
	GivenName                *string   `json:"given_name,omitempty"`
	FamilyName               *string   `json:"family_name,omitempty"`
	MiddleName               *string   `json:"middle_name,omitempty"`
	NickName                 *string   `json:"nick_name,omitempty"`
	Picture                  *string   `json:"picture,omitempty"`
	Gender                   *string   `json:"gender,omitempty"`
	BirthDate                *string   `json:"birthdate,omitempty"`
	PhoneNumber              *string   `json:"phone_number,omitempty"`
	Roles                    []*string `json:"roles,omitempty"`
	Scope                    []*string `json:"scope,omitempty"`
	RedirectURI              *string   `json:"redirect_uri,omitempty"`
	IsMultiFactorAuthEnabled *bool     `json:"is_multi_factor_auth_enabled,omitempty"`
}

// UpdateProfile is method attached to AuthorizerClient.
// It performs update_profile mutation on authorizer instance.
// It returns User reference or error.
// For implementation details check UpdateProfileExample examples/update_profile.go
func (c *AuthorizerClient) UpdateProfile(req *UpdateProfileInput, headers map[string]string) (*Response, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `mutation updateProfile($data: UpdateProfileInput!) {	update_profile(params: $data) { message } }`,
		Variables: map[string]interface{}{
			"data": req,
		},
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]*Response
	json.Unmarshal(bytesData, &res)

	return res["update_profile"], nil
}
