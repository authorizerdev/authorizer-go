package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/authorizerdev/authorizer-go"
	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
)

// These are pure unit tests (httptest, no docker server needed). Run with:
//   go test ./test/ -run 'Unit'

// tokenClient builds a user client pointed at the given httptest server.
func tokenClient(t *testing.T, serverURL string) *authorizer.AuthorizerClient {
	t.Helper()
	c, err := authorizer.NewAuthorizerClient("test-client-id", serverURL, "", map[string]string{})
	if err != nil {
		t.Fatalf("NewAuthorizerClient failed: %v", err)
	}
	return c
}

// TestGetTokenRequestConstructionUnit verifies GetToken sends a form-encoded
// POST to /oauth/token carrying exactly the set parameters (plus grant_type
// and client_id) and parses the standard token response.
func TestGetTokenRequestConstructionUnit(t *testing.T) {
	var gotForm map[string][]string
	var gotPath, gotContentType, gotOrigin string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotContentType = r.Header.Get("Content-Type")
		gotOrigin = r.Header.Get("Origin")
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm failed: %v", err)
		}
		gotForm = r.PostForm
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":      "at-123",
			"token_type":        "Bearer",
			"expires_in":        1800,
			"scope":             "read:mail",
			"issued_token_type": "urn:ietf:params:oauth:token-type:access_token",
		})
	}))
	defer srv.Close()

	c := tokenClient(t, srv.URL)
	res, err := c.GetToken(&authorizer.GetTokenRequest{
		GrantType:    authorizer.NewStringRef(authorizer.GrantTypeTokenExchange),
		SubjectToken: authorizer.NewStringRef("subj-tok"),
		ActorToken:   authorizer.NewStringRef("act-tok"),
		Scope:        authorizer.NewStringRef("read:mail"),
		Resource:     authorizer.NewStringRef("https://api.example.com"),
	})
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	if gotPath != "/oauth/token" {
		t.Errorf("expected path /oauth/token, got %q", gotPath)
	}
	if gotContentType != "application/x-www-form-urlencoded" {
		t.Errorf("expected form content type, got %q", gotContentType)
	}
	if gotOrigin == "" {
		t.Error("expected Origin header to be auto-injected (CSRF guard)")
	}
	want := map[string]string{
		"grant_type":    authorizer.GrantTypeTokenExchange,
		"client_id":     "test-client-id",
		"subject_token": "subj-tok",
		"actor_token":   "act-tok",
		"scope":         "read:mail",
		"resource":      "https://api.example.com",
	}
	for k, v := range want {
		if got := strings.Join(gotForm[k], ","); got != v {
			t.Errorf("form[%s] = %q, want %q", k, got, v)
		}
	}
	// Unset params must not be sent at all.
	for _, k := range []string{"code", "refresh_token", "client_secret", "client_assertion", "client_assertion_type", "subject_token_type", "actor_token_type"} {
		if _, ok := gotForm[k]; ok {
			t.Errorf("unset param %q was sent: %v", k, gotForm[k])
		}
	}

	if res.AccessToken != "at-123" || res.TokenType != "Bearer" || res.ExpiresIn != 1800 ||
		res.Scope != "read:mail" || res.IssuedTokenType != "urn:ietf:params:oauth:token-type:access_token" {
		t.Errorf("unexpected token response: %+v", res)
	}
}

// TestGetTokenClientCredentialsUnit verifies the client_credentials grant
// sends client_secret and parses the machine-token response.
func TestGetTokenClientCredentialsUnit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.PostForm.Get("grant_type") != authorizer.GrantTypeClientCredentials {
			t.Errorf("grant_type = %q", r.PostForm.Get("grant_type"))
		}
		if r.PostForm.Get("client_secret") != "s3cret" {
			t.Errorf("client_secret = %q", r.PostForm.Get("client_secret"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "machine-at",
			"token_type":   "Bearer",
			"expires_in":   900,
			"scope":        "a b",
		})
	}))
	defer srv.Close()

	c := tokenClient(t, srv.URL)
	res, err := c.GetToken(&authorizer.GetTokenRequest{
		GrantType:    authorizer.NewStringRef(authorizer.GrantTypeClientCredentials),
		ClientSecret: authorizer.NewStringRef("s3cret"),
		Scope:        authorizer.NewStringRef("a b"),
	})
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}
	if res.AccessToken != "machine-at" || res.Scope != "a b" {
		t.Errorf("unexpected response: %+v", res)
	}
}

// TestGetTokenErrorMappingUnit verifies RFC 6749 §5.2 error responses map to
// Go errors carrying error + error_description.
func TestGetTokenErrorMappingUnit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":             "invalid_grant",
			"error_description": "The authorization code is invalid",
		})
	}))
	defer srv.Close()

	c := tokenClient(t, srv.URL)
	_, err := c.GetToken(&authorizer.GetTokenRequest{
		Code: authorizer.NewStringRef("bad-code"),
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid_grant") || !strings.Contains(err.Error(), "The authorization code is invalid") {
		t.Errorf("expected oauth error mapping, got %v", err)
	}
}

// TestGetTokenValidationUnit verifies pre-network validation: refresh grant
// without a refresh token and default-grant (authorization_code) without a
// code both fail before any request is made.
func TestGetTokenValidationUnit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("no request should reach the server")
	}))
	defer srv.Close()
	c := tokenClient(t, srv.URL)

	if _, err := c.GetToken(&authorizer.GetTokenRequest{
		GrantType: authorizer.NewStringRef(authorizer.GrantTypeRefreshToken),
	}); err == nil {
		t.Error("expected error for refresh grant without refresh token")
	}
	if _, err := c.GetToken(&authorizer.GetTokenRequest{}); err == nil {
		t.Error("expected error for default grant without code")
	}
}

// TestAdminGraphQLWrapUnit verifies graphql ops that return a domain object
// directly (e.g. _update_client → Client) are re-wrapped onto the proto
// response envelope (UpdateClientResponse{client}).
func TestAdminGraphQLWrapUnit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/graphql" {
			t.Errorf("expected /graphql, got %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"_update_client": map[string]interface{}{
					"id":             "client-1",
					"name":           "renamed",
					"allowed_scopes": []string{"read:mail"},
					"is_active":      true,
				},
			},
		})
	}))
	defer srv.Close()

	c, err := authorizer.NewAuthorizerAdminClient(srv.URL, "admin-secret")
	if err != nil {
		t.Fatalf("NewAuthorizerAdminClient failed: %v", err)
	}
	res, err := c.UpdateClient(&authorizerv1.UpdateClientRequest{Id: "client-1"})
	if err != nil {
		t.Fatalf("UpdateClient failed: %v", err)
	}
	if res.GetClient().GetId() != "client-1" || res.GetClient().GetName() != "renamed" {
		t.Errorf("wrapped client not mapped: %+v", res)
	}
}
