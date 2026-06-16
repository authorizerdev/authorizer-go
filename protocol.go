package authorizer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Protocol selects the wire transport a client uses to talk to authorizer.
type Protocol string

const (
	// ProtocolGraphQL routes calls through POST /graphql. This is the default
	// and keeps the SDK 100% backward compatible.
	ProtocolGraphQL Protocol = "graphql"

	// ProtocolREST routes calls through the typed REST endpoints
	// (POST/GET /v1/<snake>) generated from the proto google.api.http
	// annotations.
	ProtocolREST Protocol = "rest"

	// ProtocolGRPC routes calls through the generated gRPC service stub by
	// dialing AuthorizerURL.
	ProtocolGRPC Protocol = "grpc"
)

// defaultGRPCPort is the port the authorizer server's gRPC listener binds to by
// default. It is separate from the HTTP port served from AuthorizerURL.
const defaultGRPCPort = "9091"

// grpcDial opens a gRPC client connection. When grpcEndpoint is non-empty it is
// dialed verbatim; otherwise the host is derived from authorizerURL and the gRPC
// server's default port (9091) is used, since gRPC listens on its own port and
// not the HTTP URL's port. An https:// URL (or an explicit :443 host) uses TLS;
// everything else dials insecurely, matching the typical self-hosted
// http://host:8080 deployment.
func grpcDial(authorizerURL, grpcEndpoint string) (*grpc.ClientConn, error) {
	u, err := url.Parse(authorizerURL)
	if err != nil {
		return nil, fmt.Errorf("invalid authorizerURL %q: %w", authorizerURL, err)
	}

	var creds credentials.TransportCredentials
	if u.Scheme == "https" || strings.HasSuffix(u.Host, ":443") {
		creds = credentials.NewTLS(&tls.Config{})
	} else {
		creds = insecure.NewCredentials()
	}

	host := grpcEndpoint
	if host == "" {
		// Derive the host from authorizerURL and target the gRPC port.
		host = u.Host
		if host == "" {
			// authorizerURL may be a bare host:port without a scheme.
			host = strings.TrimSuffix(authorizerURL, "/")
		}
		host = stripPort(host) + ":" + defaultGRPCPort
	}
	if (u.Scheme == "https" || strings.HasSuffix(u.Host, ":443")) && !strings.Contains(host, ":") {
		host += ":443"
	}

	return grpc.NewClient(host, grpc.WithTransportCredentials(creds))
}

// stripPort removes a trailing :port from host, leaving the bare host.
func stripPort(host string) string {
	if i := strings.LastIndex(host, ":"); i != -1 {
		return host[:i]
	}
	return host
}

// grpcContext builds a context carrying the given outgoing metadata headers.
func grpcContext(headers map[string]string) context.Context {
	return outgoingContext(context.Background(), headers)
}
