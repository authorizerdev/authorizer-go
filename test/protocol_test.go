package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/authorizerdev/authorizer-go"
)

// protocols is the set of wire transports the public client supports. Each
// integration test that should hold across transports loops over this slice.
// graphql is the default; rest and grpc are wired from the proto annotations.
var protocols = []authorizer.Protocol{
	authorizer.ProtocolGraphQL,
	authorizer.ProtocolREST,
	authorizer.ProtocolGRPC,
}

// protocolClient builds a user client pinned to the given protocol. The Origin
// header is required for the HTTP transports (CSRF guard); it is harmless for
// grpc (carried as metadata, ignored server-side).
func protocolClient(t *testing.T, p authorizer.Protocol) *authorizer.AuthorizerClient {
	t.Helper()
	c, err := authorizer.NewAuthorizerClient(clientID, authorizerURL, "", map[string]string{
		"Origin": authorizerURL,
	}, authorizer.WithProtocol(p))
	if err != nil {
		t.Fatalf("failed to create %s client: %v", p, err)
	}
	return c
}

// TestProtocolDefaultIsGraphQL guards the backward-compatibility contract: a
// client built the old way must default to graphql.
func TestProtocolDefaultIsGraphQL(t *testing.T) {
	c, err := authorizer.NewAuthorizerClient(clientID, authorizerURL, "", nil)
	if err != nil {
		t.Fatalf("NewAuthorizerClient failed: %v", err)
	}
	if c.Protocol != authorizer.ProtocolGraphQL {
		t.Errorf("expected default protocol graphql, got %q", c.Protocol)
	}
}

// TestSignUpProfileAcrossProtocols exercises the cross-protocol public flow over
// graphql+rest+grpc: signup returns an access token, and that token is used to
// fetch the profile over the same protocol. This verifies the flat-envelope
// mapping end to end: AuthResponse → AuthTokenResponse and User → User come back
// populated over all three transports (the int64 created_at/updated_at fields
// decode via protojson over REST), and that authenticated grpc Profile calls
// carry the bearer token as metadata (the #636 interceptor rejects
// unauthenticated non-public RPCs).
func TestSignUpProfileAcrossProtocols(t *testing.T) {
	for _, p := range protocols {
		p := p
		t.Run(string(p), func(t *testing.T) {
			c := protocolClient(t, p)
			email := uniqueEmail()

			signupRes, err := c.SignUp(&authorizer.SignUpRequest{
				Email:           &email,
				Password:        testPassword,
				ConfirmPassword: testPassword,
			})
			if err != nil {
				t.Fatalf("[%s] SignUp failed: %v", p, err)
			}
			if signupRes == nil || signupRes.AccessToken == nil || *signupRes.AccessToken == "" {
				t.Fatalf("[%s] SignUp: expected non-empty access_token, got %+v", p, signupRes)
			}
			if signupRes.User == nil || signupRes.User.Email != email {
				t.Fatalf("[%s] SignUp: expected user email %q, got %+v", p, email, signupRes.User)
			}

			authHeader := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", authorizer.StringValue(signupRes.AccessToken)),
			}

			profile, err := c.GetProfile(authHeader)
			if err != nil {
				t.Fatalf("[%s] GetProfile failed: %v", p, err)
			}
			if profile == nil || profile.Email != email {
				t.Errorf("[%s] GetProfile: expected email %q, got %+v", p, email, profile)
			}
		})
	}
}

// TestLoginAcrossProtocols verifies Login works over graphql, rest AND grpc as
// of 2.3.0-rc.9 (PR #635 migrated Login to the service layer). Each protocol
// must return a populated access token from the flat AuthResponse.
func TestLoginAcrossProtocols(t *testing.T) {
	for _, p := range protocols {
		p := p
		t.Run(string(p), func(t *testing.T) {
			email := uniqueEmail()
			// Signup over graphql so every protocol's Login has a fresh account.
			gql := protocolClient(t, authorizer.ProtocolGraphQL)
			if _, err := gql.SignUp(&authorizer.SignUpRequest{
				Email:           &email,
				Password:        testPassword,
				ConfirmPassword: testPassword,
			}); err != nil {
				t.Fatalf("SignUp failed: %v", err)
			}

			c := protocolClient(t, p)
			loginRes, err := c.Login(&authorizer.LoginRequest{Email: &email, Password: testPassword})
			if err != nil {
				t.Fatalf("[%s] Login failed: %v", p, err)
			}
			if loginRes == nil || loginRes.AccessToken == nil || *loginRes.AccessToken == "" {
				t.Fatalf("[%s] Login: expected non-empty access_token, got %+v", p, loginRes)
			}
		})
	}
}

// TestPublicMethodsReachServerAcrossProtocols asserts that every formerly
// graphql-only public method is now invokable over rest+grpc: the call must
// reach the server (no client-side "not available over <protocol>" error). The
// server's domain response (success or a real validation error) is accepted —
// the contract under test is protocol availability, not auth/email semantics.
func TestPublicMethodsReachServerAcrossProtocols(t *testing.T) {
	email := uniqueEmail()
	phone := "+10000000000"
	publicMethods := map[string]func(*authorizer.AuthorizerClient) error{
		"MagicLinkLogin": func(c *authorizer.AuthorizerClient) error {
			_, e := c.MagicLinkLogin(&authorizer.MagicLinkLoginRequest{Email: email})
			return e
		},
		"VerifyEmail": func(c *authorizer.AuthorizerClient) error {
			_, e := c.VerifyEmail(&authorizer.VerifyEmailRequest{Token: "t"})
			return e
		},
		"ResendVerifyEmail": func(c *authorizer.AuthorizerClient) error {
			_, e := c.ResendVerifyEmail(&authorizer.ResendVerifyEmailRequest{Email: email})
			return e
		},
		"VerifyOTP": func(c *authorizer.AuthorizerClient) error {
			_, e := c.VerifyOTP(&authorizer.VerifyOTPRequest{Email: &email, OTP: "000000"})
			return e
		},
		"ResendOTP": func(c *authorizer.AuthorizerClient) error {
			_, e := c.ResendOTP(&authorizer.ResendOTPRequest{Email: &email})
			return e
		},
		"ForgotPassword": func(c *authorizer.AuthorizerClient) error {
			_, e := c.ForgotPassword(&authorizer.ForgotPasswordRequest{Email: &email})
			return e
		},
		"ResetPassword": func(c *authorizer.AuthorizerClient) error {
			_, e := c.ResetPassword(&authorizer.ResetPasswordRequest{Password: testPassword, ConfirmPassword: testPassword, PhoneNumber: &phone})
			return e
		},
	}

	for _, p := range protocols {
		p := p
		c := protocolClient(t, p)
		for name, call := range publicMethods {
			name, call := name, call
			t.Run(fmt.Sprintf("%s/%s", p, name), func(t *testing.T) {
				// Any error must be a server/domain error, never the SDK's
				// client-side unsupported-protocol guard.
				if err := call(c); err != nil && strings.Contains(err.Error(), "not available over ") {
					t.Errorf("[%s] %s: method should be available over %s now, got %v", p, name, p, err)
				}
			})
		}
	}
}

// TestUpdateProfileAcrossProtocols verifies the authenticated UpdateProfile RPC
// works over graphql+rest+grpc with a bearer token. For grpc this exercises the
// #636 interceptor path: the bearer token must be carried as metadata or the
// non-public RPC is rejected before the handler.
func TestUpdateProfileAcrossProtocols(t *testing.T) {
	for _, p := range protocols {
		p := p
		t.Run(string(p), func(t *testing.T) {
			email := uniqueEmail()
			c := protocolClient(t, p)
			signupRes, err := c.SignUp(&authorizer.SignUpRequest{
				Email:           &email,
				Password:        testPassword,
				ConfirmPassword: testPassword,
			})
			if err != nil {
				t.Fatalf("[%s] SignUp failed: %v", p, err)
			}
			authHeader := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", authorizer.StringValue(signupRes.AccessToken)),
			}
			given := "Cross Protocol"
			res, err := c.UpdateProfile(&authorizer.UpdateProfileRequest{GivenName: &given}, authHeader)
			if err != nil {
				t.Fatalf("[%s] UpdateProfile failed: %v", p, err)
			}
			if res == nil || res.Message == "" {
				t.Errorf("[%s] UpdateProfile: expected non-empty message, got %+v", p, res)
			}
		})
	}
}

// TestGetMetaDataAcrossProtocols checks the no-auth meta endpoint over every
// protocol (graphql query, GET /v1/meta, grpc Meta).
func TestGetMetaDataAcrossProtocols(t *testing.T) {
	for _, p := range protocols {
		p := p
		t.Run(string(p), func(t *testing.T) {
			c := protocolClient(t, p)
			res, err := c.GetMetaData()
			if err != nil {
				t.Fatalf("[%s] GetMetaData failed: %v", p, err)
			}
			if res == nil || res.Version == "" {
				t.Errorf("[%s] GetMetaData: expected non-empty version, got %+v", p, res)
			}
		})
	}
}
