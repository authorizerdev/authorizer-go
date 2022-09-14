package authorizer

import "fmt"

const (
	// UserFragment defines graphql fragment for all the user attributes
	UserFragment = `
		id
		email
		email_verified
		given_name
		family_name
		middle_name
		nickname
		preferred_username
		picture
		signup_methods
		gender
		birthdate
		phone_number
		phone_number_verified
		roles
		created_at
		updated_at
		is_multi_factor_auth_enabled
	`
)

// AuthTokenResponseFragment defines graphql response for auth token type,
// which is common across various authorizer operations
var AuthTokenResponseFragment = fmt.Sprintf(`
		message
		access_token
		expires_in
		refresh_token
		id_token
		should_show_otp_screen
		user {
			%s
		}`, UserFragment,
)

// User defines attributes for user instance
type User struct {
	ID                       string    `json:"id"`
	Email                    string    `json:"email"`
	PreferredUsername        string    `json:"preferred_username"`
	EmailVerified            bool      `json:"email_verified"`
	SignupMethods            string    `json:"signup_methods"`
	GivenName                *string   `json:"given_name"`
	FamilyName               *string   `json:"family_name"`
	MiddleName               *string   `json:"middle_name"`
	Nickname                 *string   `json:"nickname"`
	Picture                  *string   `json:"picture"`
	Gender                   *string   `json:"gender"`
	Birthdate                *string   `json:"birthdate"`
	PhoneNumber              *string   `json:"phone_number"`
	PhoneNumberVerified      *bool     `json:"phone_number_verified"`
	Roles                    []*string `json:"roles"`
	CreatedAt                int64     `json:"created_at"`
	UpdatedAt                int64     `json:"updated_at"`
	IsMultiFactorAuthEnabled *bool     `json:"is_multi_factor_auth_enabled"`
}

// AuthTokenResponse defines attribute for auth token response,
// which is common across various authorizer operations
type AuthTokenResponse struct {
	Message      *string `json:"message,omitempty"`
	AccessToken  *string `json:"access_token,omitempty"`
	ExpiresIn    *int64  `json:"expires_in,omitempty"`
	IdToken      *string `json:"id_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	OtpSent      *bool   `json:"should_show_otp_screen"`
	User         *User   `json:"user,omitempty"`
}

// Response defines attribute for Response graphql type
// it is common across various authorizer operations
type Response struct {
	Message string `json:"message"`
}
