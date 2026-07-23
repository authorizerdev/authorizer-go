package authorizer

import "fmt"

const (
	// GrantTypeAuthorizationCode is used for authorization_code grant type
	GrantTypeAuthorizationCode = "authorization_code"

	// GrantTypeRefreshToken is used for refresh_token grant type
	GrantTypeRefreshToken = "refresh_token"

	// GrantTypeClientCredentials is used for the client_credentials grant type
	// (machine-to-machine tokens for registered clients / service accounts).
	GrantTypeClientCredentials = "client_credentials"

	// GrantTypeTokenExchange is the RFC 8693 token-exchange grant type
	// (delegation / impersonation via subject_token + optional actor_token).
	GrantTypeTokenExchange = "urn:ietf:params:oauth:grant-type:token-exchange"

	// UserFragment defines graphql fragment for all the user attributes
	UserFragment = `id email email_verified given_name family_name middle_name nickname preferred_username picture signup_methods gender birthdate phone_number phone_number_verified roles created_at updated_at is_multi_factor_auth_enabled has_skipped_mfa_setup_at mfa_locked_at enrolled_mfa_methods app_data revoked_timestamp
	`
)

// AuthTokenResponseFragment defines graphql response for auth token type,
// which is common across various authorizer operations
var AuthTokenResponseFragment = fmt.Sprintf(`message access_token expires_in refresh_token id_token should_show_email_otp_screen should_show_mobile_otp_screen should_show_totp_screen should_offer_webauthn_mfa_verify should_offer_webauthn_mfa_setup should_offer_email_otp_mfa_setup should_offer_sms_otp_mfa_setup authenticator_scanner_image authenticator_secret authenticator_recovery_codes user { 	%s }`, UserFragment)

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
	// HasSkippedMfaSetupAt is set once the user explicitly skips the optional
	// MFA setup prompt shown at login. Nil means never skipped.
	HasSkippedMfaSetupAt *int64 `json:"has_skipped_mfa_setup_at"`
	// MfaLockedAt is set once the user reports losing access to their only MFA
	// factor(s) with no OTP fallback enrolled. Nil means not locked.
	MfaLockedAt *int64 `json:"mfa_locked_at"`
	// EnrolledMfaMethods lists the MFA factors this user has actually
	// verified/enrolled: any of "totp", "webauthn", "email_otp", "sms_otp".
	EnrolledMfaMethods []string               `json:"enrolled_mfa_methods"`
	AppData            map[string]interface{} `json:"app_data,omitempty"`
	RevokedTimestamp   *int64                 `json:"revoked_timestamp"`
}

// AuthTokenResponse defines attribute for auth token response,
// which is common across various authorizer operations
type AuthTokenResponse struct {
	Message                   *string `json:"message,omitempty"`
	AccessToken               *string `json:"access_token,omitempty"`
	ExpiresIn                 *int64  `json:"expires_in,omitempty"`
	IdToken                   *string `json:"id_token,omitempty"`
	RefreshToken              *string `json:"refresh_token,omitempty"`
	ShouldShowEmailOtpScreen  *bool   `json:"should_show_email_otp_screen"`
	ShouldShowMobileOtpScreen *bool   `json:"should_show_mobile_otp_screen"`
	ShouldShowTotpScreen      *bool   `json:"should_show_totp_screen"`
	// ShouldOfferWebauthnMfaVerify is true when the authenticated-with-password
	// user has a registered passkey and MFA verification (not enrollment) is
	// required before a token is issued.
	ShouldOfferWebauthnMfaVerify *bool `json:"should_offer_webauthn_mfa_verify"`
	// ShouldOfferWebauthnMfaSetup / ShouldOfferEmailOtpMfaSetup /
	// ShouldOfferSmsOtpMfaSetup are true, alongside ShouldShowTotpScreen, when
	// this is a first-time optional-MFA offer and that method is available.
	// AccessToken is NOT populated alongside these flags; it is withheld until
	// the user completes a method or skips (see SkipMfaSetup).
	ShouldOfferWebauthnMfaSetup *bool     `json:"should_offer_webauthn_mfa_setup"`
	ShouldOfferEmailOtpMfaSetup *bool     `json:"should_offer_email_otp_mfa_setup"`
	ShouldOfferSmsOtpMfaSetup   *bool     `json:"should_offer_sms_otp_mfa_setup"`
	AuthenticatorScannerImage   *string   `json:"authenticator_scanner_image"`
	AuthenticatorSecret         *string   `json:"authenticator_secret"`
	AuthenticatorRecoveryCodes  []*string `json:"authenticator_recovery_codes"`
	User                        *User     `json:"user,omitempty"`
}

// Response defines attribute for Response graphql type
// it is common across various authorizer operations
type Response struct {
	Message string `json:"message"`
}

// ForgotPasswordResponse defines attribute for forgot_password response
type ForgotPasswordResponse struct {
	Message                   string `json:"message"`
	ShouldShowMobileOtpScreen *bool  `json:"should_show_mobile_otp_screen"`
}
