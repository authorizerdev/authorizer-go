// Package authorizer provides client and methods using which you can
// perform various graphql operations to your authorizer instance
package authorizer

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc"
)

// AuthorizerClient defines the attributes required to initiate authorizer client
type AuthorizerClient struct {
	ClientID      string
	AuthorizerURL string
	RedirectURL   string
	ExtraHeaders  map[string]string
	// Protocol selects the wire transport (graphql, rest or grpc). Defaults to
	// ProtocolGraphQL when unset, keeping the SDK backward compatible.
	Protocol Protocol
	// GRPCEndpoint overrides the host:port dialed when Protocol is grpc. When
	// empty it is derived from AuthorizerURL using the gRPC default port.
	GRPCEndpoint string

	// httpClient is shared across every call on this client and carries a
	// cookie jar. The jar is REQUIRED, not an optimisation: the MFA offer flow
	// identifies the pending user by a session cookie the server sets on
	// signup/login, and SkipMfaSetup / VerifyOtp only resolve it if that cookie
	// is sent back. Building a fresh http.Client per call, as this SDK used to,
	// dropped the cookie and made those calls fail with "invalid session" —
	// i.e. the whole MFA surface was unreachable from Go.
	httpClient *http.Client
}

// newHTTPClient builds the shared cookie-aware client. cookiejar.New with a nil
// options value never returns an error, but the error is handled rather than
// ignored so a future options change cannot silently produce a jar-less client.
func newHTTPClient() *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &http.Client{Timeout: 30 * time.Second}
	}
	return &http.Client{Timeout: 30 * time.Second, Jar: jar}
}

// HTTPClient returns the shared cookie-aware http client, initialising it on
// first use so a zero-value AuthorizerClient (or one built by an older
// constructor path) still carries a jar.
func (c *AuthorizerClient) HTTPClient() *http.Client {
	if c.httpClient == nil {
		c.httpClient = newHTTPClient()
	}
	return c.httpClient
}

// dialGRPC opens a gRPC connection whose calls share this client's cookie jar,
// so a session established over gRPC (or over HTTP) survives the next call.
func (c *AuthorizerClient) dialGRPC() (*grpc.ClientConn, error) {
	jar := c.HTTPClient().Jar
	u, err := url.Parse(c.AuthorizerURL)
	if jar == nil || err != nil || u.Host == "" {
		return grpcDial(c.AuthorizerURL, c.GRPCEndpoint)
	}
	return grpcDial(c.AuthorizerURL, c.GRPCEndpoint,
		grpc.WithChainUnaryInterceptor(cookieInterceptor(jar, u)))
}

// ClientOption customizes an AuthorizerClient at construction time.
type ClientOption func(*AuthorizerClient)

// WithProtocol sets the wire transport the client uses (graphql, rest or grpc).
func WithProtocol(p Protocol) ClientOption {
	return func(c *AuthorizerClient) {
		c.Protocol = p
	}
}

// WithGRPCEndpoint sets the host:port dialed for grpc calls. The authorizer
// server's gRPC listener runs on its own port (default 9091), separate from the
// HTTP port in AuthorizerURL. When unset, the endpoint is derived from
// AuthorizerURL's host with the default gRPC port (9091).
func WithGRPCEndpoint(addr string) ClientOption {
	return func(c *AuthorizerClient) {
		c.GRPCEndpoint = addr
	}
}

// NewAuthorizerClient creates an authorizer client instance.
// It returns reference to authorizer client instance or error.
// The optional functional options (e.g. WithProtocol) tweak behavior while
// keeping the original positional signature backward compatible.
func NewAuthorizerClient(clientID, authorizerURL, redirectURL string, extraHeaders map[string]string, opts ...ClientOption) (*AuthorizerClient, error) {
	if strings.TrimSpace(clientID) == "" {
		return nil, fmt.Errorf("clientID missing")
	}

	if strings.TrimSpace(authorizerURL) == "" {
		return nil, fmt.Errorf("authorizerURL missing")
	}

	// extraHeaders is optional parameter,
	// hence if not set, initialize it with empty map
	headers := extraHeaders
	if headers == nil {
		headers = make(map[string]string)
	}

	// if x-authorizer-url is not present
	// set it to authorizerURL
	if _, ok := headers["x-authorizer-url"]; !ok {
		headers["x-authorizer-url"] = authorizerURL
	}

	// Add clientID to headers
	headers["x-authorizer-client-id"] = clientID

	c := &AuthorizerClient{
		RedirectURL:   strings.TrimSuffix(redirectURL, "/"),
		AuthorizerURL: strings.TrimSuffix(authorizerURL, "/"),
		ClientID:      clientID,
		ExtraHeaders:  headers,
		Protocol:      ProtocolGraphQL,
		httpClient:    newHTTPClient(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}
