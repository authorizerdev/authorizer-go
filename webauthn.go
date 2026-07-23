package authorizer

// WebAuthn / passkey self-service ceremonies. GraphQL only: the server has no
// gRPC/REST RPCs for these operations. `options` and `credential` are opaque
// JSON strings carrying the WebAuthn PublicKeyCredential* structures; the SDK
// caller performs the base64url <-> ArrayBuffer conversion between these
// strings and the browser's navigator.credentials API.

// WebauthnCredentialInfo defines attributes for a registered passkey.
type WebauthnCredentialInfo struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Transports []*string `json:"transports"`
	CreatedAt  *int64    `json:"created_at"`
	UpdatedAt  *int64    `json:"updated_at"`
	LastUsedAt *int64    `json:"last_used_at"`
}

// WebauthnRegistrationOptionsRequest defines attributes for
// webauthn_registration_options. Email/PhoneNumber are only used in the
// MFA-session-cookie mode (a caller in the withheld first-time-offer state,
// with no bearer token yet); ignored for an ordinary bearer-token-
// authenticated settings-page caller.
type WebauthnRegistrationOptionsRequest struct {
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

// WebauthnRegistrationOptionsResponse defines attributes for
// webauthn_registration_options response.
type WebauthnRegistrationOptionsResponse struct {
	// Options is JSON-encoded PublicKeyCredentialCreationOptions to pass to
	// navigator.credentials.create().
	Options string `json:"options"`
}

// WebauthnRegistrationOptions is method attached to AuthorizerClient.
// It performs webauthn_registration_options mutation on authorizer instance,
// returning the passkey creation challenge.
func (c *AuthorizerClient) WebauthnRegistrationOptions(req *WebauthnRegistrationOptionsRequest, headers map[string]string) (*WebauthnRegistrationOptionsResponse, error) {
	var res WebauthnRegistrationOptionsResponse
	err := c.execute(methodSpec{
		name: "WebauthnRegistrationOptions",
		graphql: &GraphQLRequest{
			Query:     `mutation webauthnRegistrationOptions($email: String, $phone_number: String) { webauthn_registration_options(email: $email, phone_number: $phone_number) { options } }`,
			Variables: map[string]interface{}{"email": req.Email, "phone_number": req.PhoneNumber},
		},
		graphqlField: "webauthn_registration_options",
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// WebauthnRegistrationVerifyRequest defines attributes for
// webauthn_registration_verify request. Email/PhoneNumber/State are only used
// on the MFA-session-cookie path; ignored for an ordinary bearer-token-
// authenticated settings-page caller.
type WebauthnRegistrationVerifyRequest struct {
	// Name is an optional human label for the passkey (e.g. "MacBook Touch ID").
	Name *string `json:"name,omitempty"`
	// Credential is the JSON-encoded PublicKeyCredential attestation response
	// from navigator.credentials.create().
	Credential  string  `json:"credential"`
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	State       *string `json:"state,omitempty"`
}

// WebauthnRegistrationVerify is method attached to AuthorizerClient.
// It performs webauthn_registration_verify mutation on authorizer instance,
// completing passkey enrollment. Returns AuthTokenResponse with an
// access_token set only on the MFA-session-cookie path (nil, message-only,
// for the ordinary authenticated-settings-page caller).
func (c *AuthorizerClient) WebauthnRegistrationVerify(req *WebauthnRegistrationVerifyRequest, headers map[string]string) (*AuthTokenResponse, error) {
	var res AuthTokenResponse
	err := c.execute(methodSpec{
		name: "WebauthnRegistrationVerify",
		graphql: &GraphQLRequest{
			Query:     `mutation webauthnRegistrationVerify($data: WebauthnRegistrationVerifyRequest!) { webauthn_registration_verify(params: $data) { ` + AuthTokenResponseFragment + ` } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "webauthn_registration_verify",
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// WebauthnLoginOptionsResponse defines attributes for webauthn_login_options
// response.
type WebauthnLoginOptionsResponse struct {
	// Options is JSON-encoded PublicKeyCredentialRequestOptions to pass to
	// navigator.credentials.get().
	Options string `json:"options"`
}

// WebauthnLoginOptions is method attached to AuthorizerClient.
// It performs webauthn_login_options mutation on authorizer instance,
// returning the passkey assertion challenge. Email is optional: omit for
// discoverable-credential (usernameless) login.
func (c *AuthorizerClient) WebauthnLoginOptions(email *string) (*WebauthnLoginOptionsResponse, error) {
	var res WebauthnLoginOptionsResponse
	err := c.execute(methodSpec{
		name: "WebauthnLoginOptions",
		graphql: &GraphQLRequest{
			Query:     `mutation webauthnLoginOptions($email: String) { webauthn_login_options(email: $email) { options } }`,
			Variables: map[string]interface{}{"email": email},
		},
		graphqlField: "webauthn_login_options",
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// WebauthnLoginVerifyRequest defines attributes for webauthn_login_verify
// request.
type WebauthnLoginVerifyRequest struct {
	// State is the OAuth authorize state to continue an in-progress
	// authorization, if any.
	State *string `json:"state,omitempty"`
	// Credential is the JSON-encoded PublicKeyCredential assertion response
	// from navigator.credentials.get().
	Credential string `json:"credential"`
}

// WebauthnLoginVerify is method attached to AuthorizerClient.
// It performs webauthn_login_verify mutation on authorizer instance,
// completing passkey login.
// It returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) WebauthnLoginVerify(req *WebauthnLoginVerifyRequest) (*AuthTokenResponse, error) {
	var res AuthTokenResponse
	err := c.execute(methodSpec{
		name: "WebauthnLoginVerify",
		graphql: &GraphQLRequest{
			Query:     `mutation webauthnLoginVerify($data: WebauthnLoginVerifyRequest!) { webauthn_login_verify(params: $data) { ` + AuthTokenResponseFragment + ` } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "webauthn_login_verify",
	}, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// WebauthnDeleteCredential is method attached to AuthorizerClient.
// It performs webauthn_delete_credential mutation on authorizer instance,
// removing one of the authenticated caller's own registered passkeys.
// DESTRUCTIVE: the passkey can no longer be used to log in.
func (c *AuthorizerClient) WebauthnDeleteCredential(id string, headers map[string]string) (*Response, error) {
	var res Response
	err := c.execute(methodSpec{
		name: "WebauthnDeleteCredential",
		graphql: &GraphQLRequest{
			Query:     `mutation webauthnDeleteCredential($id: ID!) { webauthn_delete_credential(id: $id) { message } }`,
			Variables: map[string]interface{}{"id": id},
		},
		graphqlField: "webauthn_delete_credential",
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// WebauthnCredentials is method attached to AuthorizerClient.
// It performs the webauthn_credentials query on authorizer instance, listing
// the authenticated caller's own registered passkeys.
func (c *AuthorizerClient) WebauthnCredentials(headers map[string]string) ([]*WebauthnCredentialInfo, error) {
	var res []*WebauthnCredentialInfo
	err := c.execute(methodSpec{
		name: "WebauthnCredentials",
		graphql: &GraphQLRequest{
			Query: `query webauthnCredentials { webauthn_credentials { id name transports created_at updated_at last_used_at } }`,
		},
		graphqlField: "webauthn_credentials",
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
