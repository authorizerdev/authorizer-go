package authorizer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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
func grpcDial(authorizerURL, grpcEndpoint string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
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

	return grpc.NewClient(host, append([]grpc.DialOption{grpc.WithTransportCredentials(creds)}, opts...)...)
}

// cookieInterceptor carries the session cookie across gRPC calls. gRPC has no
// cookie concept, so the server sends its cookies as `set-cookie` header
// metadata and reads them back from a `cookie` metadata entry. Without this the
// MFA offer flow is unreachable over gRPC for the same reason it was over
// HTTP before the jar existed: signup/login hand out an MFA session the next
// call never replays, and SkipMfaSetup/VerifyOtp answer "invalid session".
// Cookies are stored in the client's shared jar, so a session started over one
// protocol is usable from another.
func cookieInterceptor(jar http.CookieJar, u *url.URL) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if cookies := jar.Cookies(u); len(cookies) > 0 {
			pairs := make([]string, 0, len(cookies))
			for _, c := range cookies {
				pairs = append(pairs, c.Name+"="+c.Value)
			}
			ctx = metadata.AppendToOutgoingContext(ctx, "cookie", strings.Join(pairs, "; "))
		}

		var header metadata.MD
		err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Header(&header))...)
		// Store cookies even when the call failed: a rejected MFA attempt can
		// still rotate the session.
		if set := header.Get("set-cookie"); len(set) > 0 {
			jar.SetCookies(u, (&http.Response{Header: http.Header{"Set-Cookie": set}}).Cookies())
		}
		return err
	}
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
