package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/authorizerdev/authorizer-go"
)

const (
	authorizerURL = "http://localhost:8080"
	clientID      = "123456"
	testPassword  = "Abc@123"
)

// testClient returns a new authorizer client configured for integration tests
func testClient(t *testing.T) *authorizer.AuthorizerClient {
	t.Helper()
	c, err := authorizer.NewAuthorizerClient(clientID, authorizerURL, "", nil)
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
