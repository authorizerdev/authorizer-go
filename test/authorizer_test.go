package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/authorizerdev/authorizer-go"
)

// Integration-test target. Overridable via env so the same suite runs against a
// local container on a non-default port and in CI.
var (
	authorizerURL = envOr("AUTHORIZER_TEST_URL", "http://localhost:8080")
	clientID      = envOr("AUTHORIZER_TEST_CLIENT_ID", "123456")
)

const testPassword = "Abc@123"

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// testClient returns a new authorizer client configured for integration tests.
// The Origin header is required: the server's CSRF middleware rejects
// state-changing requests (all GraphQL POSTs) that carry neither an Origin nor
// a Referer header, and in wildcard allowed-origins mode the origin must match
// the server's own host.
func testClient(t *testing.T) *authorizer.AuthorizerClient {
	t.Helper()
	c, err := authorizer.NewAuthorizerClient(clientID, authorizerURL, "", map[string]string{
		"Origin": authorizerURL,
	})
	if err != nil {
		t.Fatalf("failed to create authorizer client: %v", err)
	}
	return c
}

// uniqueEmail generates a unique email for each test run to avoid conflicts
func uniqueEmail() string {
	return fmt.Sprintf("test-%d@yopmail.com", rand.Int63())
}

func TestGetMetaData(t *testing.T) {
	c := testClient(t)

	res, err := c.GetMetaData()
	if err != nil {
		t.Fatalf("GetMetaData failed: %v", err)
	}

	if res == nil {
		t.Fatal("GetMetaData returned nil response")
	}
	if res.ClientID == "" {
		t.Error("GetMetaData: expected non-empty client_id")
	}
	if res.Version == "" {
		t.Error("GetMetaData: expected non-empty version")
	}
}

func TestSignUp(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()

	res, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           &email,
		Password:        testPassword,
		ConfirmPassword: testPassword,
	})
	if err != nil {
		t.Fatalf("SignUp failed: %v", err)
	}

	if res == nil {
		t.Fatal("SignUp returned nil response")
	}
	if res.Message != nil && *res.Message != "" {
		t.Logf("SignUp message: %s", *res.Message)
	}
}

func TestLogin(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()

	// Sign up first to create a user
	_, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           &email,
		Password:        testPassword,
		ConfirmPassword: testPassword,
	})
	if err != nil {
		t.Fatalf("SignUp failed (prerequisite for Login): %v", err)
	}

	res, err := c.Login(&authorizer.LoginRequest{
		Email:    &email,
		Password: testPassword,
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if res == nil {
		t.Fatal("Login returned nil response")
	}
	if res.AccessToken == nil || *res.AccessToken == "" {
		t.Error("Login: expected non-empty access_token")
	}
}

func TestGetProfile(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()

	// Sign up and login first
	_, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           &email,
		Password:        testPassword,
		ConfirmPassword: testPassword,
	})
	if err != nil {
		t.Fatalf("SignUp failed (prerequisite): %v", err)
	}

	loginRes, err := c.Login(&authorizer.LoginRequest{
		Email:    &email,
		Password: testPassword,
	})
	if err != nil {
		t.Fatalf("Login failed (prerequisite): %v", err)
	}

	res, err := c.GetProfile(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authorizer.StringValue(loginRes.AccessToken)),
	})
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}

	if res == nil {
		t.Fatal("GetProfile returned nil response")
	}
	if res.Email != email {
		t.Errorf("GetProfile: expected email %q, got %q", email, res.Email)
	}
}

func TestGetSession(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()

	_, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           &email,
		Password:        testPassword,
		ConfirmPassword: testPassword,
	})
	if err != nil {
		t.Fatalf("SignUp failed (prerequisite): %v", err)
	}

	loginRes, err := c.Login(&authorizer.LoginRequest{
		Email:    &email,
		Password: testPassword,
	})
	if err != nil {
		t.Fatalf("Login failed (prerequisite): %v", err)
	}

	res, err := c.GetSession(&authorizer.SessionQueryRequest{
		Roles: []*string{},
	}, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authorizer.StringValue(loginRes.AccessToken)),
	})
	if err != nil {
		// Session endpoint may return unauthorized depending on authorizer config
		if err.Error() == "unauthorized" {
			t.Skip("GetSession returned unauthorized - session API may require additional authorizer configuration")
		}
		t.Fatalf("GetSession failed: %v", err)
	}

	if res == nil {
		t.Fatal("GetSession returned nil response")
	}
}

func TestLogout(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()

	_, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           &email,
		Password:        testPassword,
		ConfirmPassword: testPassword,
	})
	if err != nil {
		t.Fatalf("SignUp failed (prerequisite): %v", err)
	}

	loginRes, err := c.Login(&authorizer.LoginRequest{
		Email:    &email,
		Password: testPassword,
	})
	if err != nil {
		t.Fatalf("Login failed (prerequisite): %v", err)
	}

	res, err := c.Logout(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authorizer.StringValue(loginRes.AccessToken)),
	})
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	if res == nil {
		t.Fatal("Logout returned nil response")
	}
	if res.Message == "" {
		t.Error("Logout: expected non-empty message")
	}
}

func TestValidateJWTToken(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()

	_, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           &email,
		Password:        testPassword,
		ConfirmPassword: testPassword,
	})
	if err != nil {
		t.Fatalf("SignUp failed (prerequisite): %v", err)
	}

	loginRes, err := c.Login(&authorizer.LoginRequest{
		Email:    &email,
		Password: testPassword,
	})
	if err != nil {
		t.Fatalf("Login failed (prerequisite): %v", err)
	}

	res, err := c.ValidateJWTToken(&authorizer.ValidateJWTTokenRequest{
		TokenType: authorizer.TokenTypeAccessToken,
		Token:     authorizer.StringValue(loginRes.AccessToken),
	})
	if err != nil {
		t.Fatalf("ValidateJWTToken failed: %v", err)
	}

	if res == nil {
		t.Fatal("ValidateJWTToken returned nil response")
	}
	if !res.IsValid {
		t.Error("ValidateJWTToken: expected token to be valid")
	}
}

func TestForgotPassword(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()

	// Create user first
	_, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           &email,
		Password:        testPassword,
		ConfirmPassword: testPassword,
	})
	if err != nil {
		t.Fatalf("SignUp failed (prerequisite): %v", err)
	}

	res, err := c.ForgotPassword(&authorizer.ForgotPasswordRequest{
		Email: &email,
	})
	if err != nil {
		t.Fatalf("ForgotPassword failed: %v", err)
	}

	if res == nil {
		t.Fatal("ForgotPassword returned nil response")
	}
	if res.Message == "" {
		t.Error("ForgotPassword: expected non-empty message")
	}
}

func TestMagicLinkLogin(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()

	// Create user first
	_, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           &email,
		Password:        testPassword,
		ConfirmPassword: testPassword,
	})
	if err != nil {
		t.Fatalf("SignUp failed (prerequisite): %v", err)
	}

	res, err := c.MagicLinkLogin(&authorizer.MagicLinkLoginRequest{
		Email: email,
	})
	if err != nil {
		// Magic link may be disabled in default authorizer config
		if err.Error() == "magic link login is disabled for this instance" {
			t.Skip("Magic link login is disabled - enable with --enable-magic-link in authorizer config")
		}
		t.Fatalf("MagicLinkLogin failed: %v", err)
	}

	if res == nil {
		t.Fatal("MagicLinkLogin returned nil response")
	}
	if res.Message == "" {
		t.Error("MagicLinkLogin: expected non-empty message")
	}
}

// skipIfFgaUnavailable skips an FGA integration test when the target server has
// the fine-grained authorization engine disabled or no model installed, so the
// suite stays green on a default (auth-only) deployment.
func skipIfFgaUnavailable(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	// The server keeps engine errors opaque: "fine-grained authorization is
	// not enabled" when started without an FGA store, and "authorization
	// check failed" / "authorization list failed" when the engine is up but
	// no authorization model has been written yet. Servers that predate the
	// FGA schema reject the query itself with `Unknown type
	// "CheckPermissionsInput"` (or the List equivalent).
	msg := err.Error()
	for _, s := range []string{"not enabled", "unauthorized", "check failed", "list failed", "unknown type"} {
		if strings.Contains(strings.ToLower(msg), s) {
			t.Skipf("FGA not available on target server (%v) - skipping", err)
		}
	}
	t.Fatalf("FGA call failed: %v", err)
}

func TestCheckPermissions(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()

	_, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           &email,
		Password:        testPassword,
		ConfirmPassword: testPassword,
	})
	if err != nil {
		t.Fatalf("SignUp failed (prerequisite): %v", err)
	}

	loginRes, err := c.Login(&authorizer.LoginRequest{
		Email:    &email,
		Password: testPassword,
	})
	if err != nil {
		t.Fatalf("Login failed (prerequisite): %v", err)
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authorizer.StringValue(loginRes.AccessToken)),
	}

	// A freshly signed up user has no relationship tuples, so every check for
	// an arbitrary object must come back denied (never errors on a healthy FGA
	// deployment). Results echo each relation/object pair in request order.
	res, err := c.CheckPermissions(&authorizer.CheckPermissionsRequest{
		Checks: []*authorizer.PermissionCheckInput{
			{Relation: "can_view", Object: "document:1"},
			{Relation: "can_edit", Object: "document:1"},
		},
	}, headers)
	skipIfFgaUnavailable(t, err)

	if res == nil || len(res.Results) != 2 {
		t.Fatalf("CheckPermissions: expected 2 results, got %+v", res)
	}
	for i, r := range res.Results {
		if r.Allowed {
			t.Errorf("CheckPermissions: expected check %d to be denied for a new user", i)
		}
		if r.Relation == "" || r.Object == "" {
			t.Errorf("CheckPermissions: expected result %d to echo the checked relation/object pair, got %+v", i, r)
		}
	}
	if res.Results[0].Relation != "can_view" || res.Results[1].Relation != "can_edit" {
		t.Errorf("CheckPermissions: results not positionally aligned with request: %+v", res.Results)
	}
}

func TestListPermissions(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()

	_, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           &email,
		Password:        testPassword,
		ConfirmPassword: testPassword,
	})
	if err != nil {
		t.Fatalf("SignUp failed (prerequisite): %v", err)
	}

	loginRes, err := c.Login(&authorizer.LoginRequest{
		Email:    &email,
		Password: testPassword,
	})
	if err != nil {
		t.Fatalf("Login failed (prerequisite): %v", err)
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authorizer.StringValue(loginRes.AccessToken)),
	}

	// A new user relates to no documents.
	objs, err := c.ListPermissions(&authorizer.ListPermissionsRequest{
		Relation:   "can_view",
		ObjectType: "document",
	}, headers)
	skipIfFgaUnavailable(t, err)
	if objs == nil {
		t.Fatal("ListPermissions returned nil response")
	}
	if len(objs.Objects) != 0 {
		t.Errorf("ListPermissions: expected no accessible objects for a new user, got %d", len(objs.Objects))
	}
	if len(objs.Permissions) != 0 {
		t.Errorf("ListPermissions: expected no permissions for a new user, got %d", len(objs.Permissions))
	}
	if objs.Truncated {
		t.Error("ListPermissions: expected truncated=false for an empty result")
	}
}
