package authorizer

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

// TestRESTTransportPersistsCookies pins the same contract for the REST
// transport, which kept using http.DefaultClient — jar-less — after the jar was
// added, so MFA over rest failed with "invalid session" while graphql worked.
func TestRESTTransportPersistsCookies(t *testing.T) {
	var secondRequestCookie string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie("mfa_session"); err == nil {
			secondRequestCookie = c.Value
		}
		http.SetCookie(w, &http.Cookie{Name: "mfa_session", Value: "session-abc", Path: "/"})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c, err := NewAuthorizerClient("test-client", srv.URL, "", nil, WithProtocol(ProtocolREST))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	for i := 0; i < 2; i++ {
		if err := c.executeREST(http.MethodPost, "/v1/signup", map[string]string{}, nil, nil); err != nil {
			t.Fatalf("rest call %d: %v", i, err)
		}
	}

	if secondRequestCookie != "session-abc" {
		t.Fatalf("cookie not replayed on the second REST call (got %q)", secondRequestCookie)
	}
}

// TestGRPCCookieInterceptor pins the gRPC half of the same contract: the
// server hands cookies out as `set-cookie` header metadata and expects them
// back as a `cookie` entry, so without this the MFA flow is unreachable over
// gRPC. The invoker is stubbed — the contract under test is metadata handling,
// not the wire.
func TestGRPCCookieInterceptor(t *testing.T) {
	u, _ := url.Parse("http://authorizer.test")
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("jar: %v", err)
	}
	intercept := cookieInterceptor(jar, u)

	// First call: the server sets a cookie, the client sends none.
	var sentFirst []string
	err = intercept(context.Background(), "/Signup", nil, nil, nil,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			md, _ := metadata.FromOutgoingContext(ctx)
			sentFirst = md.Get("cookie")
			setHeader(opts, metadata.Pairs("set-cookie", "mfa_session=session-abc; Path=/"))
			return nil
		})
	if err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(sentFirst) != 0 {
		t.Errorf("expected no cookie metadata on the first call, got %v", sentFirst)
	}

	// Second call: the stored cookie must be replayed.
	var sentSecond []string
	err = intercept(context.Background(), "/SkipMfaSetup", nil, nil, nil,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			md, _ := metadata.FromOutgoingContext(ctx)
			sentSecond = md.Get("cookie")
			return nil
		})
	if err != nil {
		t.Fatalf("second call: %v", err)
	}
	if len(sentSecond) != 1 || !strings.Contains(sentSecond[0], "mfa_session=session-abc") {
		t.Fatalf("cookie not replayed over grpc (got %v) — SkipMfaSetup will fail with \"invalid session\"", sentSecond)
	}
}

// setHeader writes response header metadata the way a real grpc call does, by
// filling the grpc.Header(&md) call option the interceptor appended.
func setHeader(opts []grpc.CallOption, md metadata.MD) {
	for _, o := range opts {
		if h, ok := o.(grpc.HeaderCallOption); ok {
			*h.HeaderAddr = md
		}
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
