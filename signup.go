package authorizer

import (
	"encoding/json"
	"fmt"
)

// SignUpRequest defines attributes for signup request
type SignUpRequest struct {
	Email                    *string                `json:"email,omitempty"`
	Password                 string                 `json:"password"`
	ConfirmPassword          string                 `json:"confirm_password"`
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
	State                    *string                `json:"state,omitempty"`
}

// SignUpInput is deprecated: Use SignUpRequest instead
type SignUpInput = SignUpRequest

// SignUp is method attached to AuthorizerClient.
// It performs signup mutation on authorizer instance.
// It takes SignUpRequest reference as parameter and returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) SignUp(req *SignUpRequest) (*AuthTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`mutation signup($data: SignUpRequest!) { signup(params: $data) { %s }}`, AuthTokenResponseFragment),
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*AuthTokenResponse
	json.Unmarshal(bytesData, &res)

	return res["signup"], nil
}
