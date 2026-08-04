package authorizer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestClientPersistsCookiesAcrossCalls is the regression test for the MFA
// offer flow.
//
// Since server 2.4.0, MFA is on by default: signup/login withhold the access
// token and return "Proceed to mfa setup", identifying the pending user by a
// session cookie. SkipMfaSetup and VerifyOtp resolve that user ONLY if the
// cookie is sent back.
//
// This SDK previously built a fresh http.Client per call, so the cookie was
// dropped between them and every one of those calls failed with "invalid
// session" — the entire MFA surface was unreachable from Go, while the methods
// existed and looked correct.
func TestClientPersistsCookiesAcrossCalls(t *testing.T) {
	var secondRequestCookie string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie("mfa_session"); err == nil {
			secondRequestCookie = c.Value
		}
		http.SetCookie(w, &http.Cookie{Name: "mfa_session", Value: "session-abc", Path: "/"})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
	}))
	defer srv.Close()

	c, err := NewAuthorizerClient("test-client", srv.URL, "", nil)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	// First call receives the cookie.
	if _, err := c.ExecuteGraphQL(&GraphQLRequest{Query: "{ __typename }"}, nil); err != nil {
		t.Fatalf("first call: %v", err)
	}
	// Second call must send it back.
	if _, err := c.ExecuteGraphQL(&GraphQLRequest{Query: "{ __typename }"}, nil); err != nil {
		t.Fatalf("second call: %v", err)
	}

	if secondRequestCookie != "session-abc" {
		t.Fatalf("cookie not replayed on the second call (got %q) — "+
			"SkipMfaSetup/VerifyOtp will fail with \"invalid session\"", secondRequestCookie)
	}
}

// TestHTTPClientAlwaysHasAJar pins that no construction path yields a
// jar-less client, including a zero-value struct built without the constructor.
func TestHTTPClientAlwaysHasAJar(t *testing.T) {
	c, err := NewAuthorizerClient("test-client", "http://localhost:8080", "", nil)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	if c.HTTPClient().Jar == nil {
		t.Fatal("constructor produced a client with no cookie jar")
	}

	var zero AuthorizerClient
	if zero.HTTPClient().Jar == nil {
		t.Fatal("zero-value client must lazily gain a cookie jar")
	}
}
