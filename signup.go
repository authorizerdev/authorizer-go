package authorizer

import (
	"encoding/json"
	"fmt"
)

// SignUpInput defines attributes for signup request
type SignUpInput struct {
	Email                    string    `json:"email"`
	Password                 string    `json:"password"`
	ConfirmPassword          string    `json:"confirm_password"`
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
	AppData                  *string   `json:"app_data,omitempty"`
}

// SignUp is method attached to AuthorizerClient.
// It performs signup mutation on authorizer instance.
// It takes SignUpInput reference as parameter and returns AuthTokenResponse reference or error.
// For implementation details check SignUpExample examples/signup.go
func (c *AuthorizerClient) SignUp(req *SignUpInput) (*AuthTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`mutation signup($data: SignUpInput!) { signup(params: $data) { %s }}`, AuthTokenResponseFragment),
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
