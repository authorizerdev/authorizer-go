// Package authorizer provides client and methods using which you can
// perform various graphql operations to your authorizer instance
package authorizer

import (
	"fmt"
	"strings"
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
		ExtraHeaders:  extraHeaders,
		Protocol:      ProtocolGraphQL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}
