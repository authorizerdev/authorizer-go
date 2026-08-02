package authorizer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	authorizerv1 "github.com/authorizerdev/authorizer-proto-go/authorizer/v1"
)

// adminSecretHeader is the header (and gRPC metadata key) carrying the admin
// secret on every admin call.
const adminSecretHeader = "x-authorizer-admin-secret"

// AuthorizerAdminClient is the admin surface of the SDK. It mirrors the
// AuthorizerAdminService proto and carries the admin secret plus the selected
// wire protocol. The request/response types are the generated proto messages
// (authorizerv1.*), the canonical signatures shared across all three SDKs.
type AuthorizerAdminClient struct {
	AuthorizerURL string
	AdminSecret   string
	ExtraHeaders  map[string]string
	// Protocol selects the wire transport. Defaults to ProtocolGraphQL.
	Protocol Protocol
	// GRPCEndpoint overrides the host:port dialed when Protocol is grpc. When
	// empty it is derived from AuthorizerURL using the gRPC default port.
	GRPCEndpoint string
}

// AdminClientOption customizes an AuthorizerAdminClient at construction time.
type AdminClientOption func(*AuthorizerAdminClient)

// WithAdminProtocol sets the wire transport the admin client uses.
func WithAdminProtocol(p Protocol) AdminClientOption {
	return func(c *AuthorizerAdminClient) {
		c.Protocol = p
	}
}

// WithAdminGRPCEndpoint sets the host:port dialed for grpc calls. The authorizer
// server's gRPC listener runs on its own port (default 9091), separate from the
// HTTP port in AuthorizerURL. When unset, the endpoint is derived from
// AuthorizerURL's host with the default gRPC port (9091).
func WithAdminGRPCEndpoint(addr string) AdminClientOption {
	return func(c *AuthorizerAdminClient) {
		c.GRPCEndpoint = addr
	}
}

// WithAdminExtraHeaders sets default headers added to every admin HTTP request.
func WithAdminExtraHeaders(headers map[string]string) AdminClientOption {
	return func(c *AuthorizerAdminClient) {
		c.ExtraHeaders = headers
	}
}

// NewAuthorizerAdminClient creates an admin client instance authenticated with
// the admin secret. Default protocol is graphql; override with WithAdminProtocol.
func NewAuthorizerAdminClient(authorizerURL, adminSecret string, opts ...AdminClientOption) (*AuthorizerAdminClient, error) {
	if strings.TrimSpace(authorizerURL) == "" {
		return nil, fmt.Errorf("authorizerURL missing")
	}
	if strings.TrimSpace(adminSecret) == "" {
		return nil, fmt.Errorf("adminSecret missing")
	}

	c := &AuthorizerAdminClient{
		AuthorizerURL: strings.TrimSuffix(authorizerURL, "/"),
		AdminSecret:   adminSecret,
		ExtraHeaders:  map[string]string{},
		Protocol:      ProtocolGraphQL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// adminMethodSpec declares how one admin method is carried over each protocol.
// A protocol whose field is empty/nil is unsupported and yields a clear error.
type adminMethodSpec struct {
	name string // human method name for error messages, e.g. "AdminMeta"

	// graphql is the prebuilt GraphQL request; nil means gql-unsupported.
	graphql *GraphQLRequest
	// graphqlField is the top-level response field to unwrap (the _-prefixed op).
	graphqlField string
	// graphqlWrap, when non-empty, re-wraps the unwrapped graphql payload as
	// {graphqlWrap: <payload>} before unmarshalling into out. Used when the
	// GraphQL op returns a domain object directly but the proto response wraps
	// it (e.g. _update_client returns Client, proto UpdateClientResponse{client}).
	graphqlWrap string

	// responseUnwrap is the dual of graphqlWrap, for methods whose out is the
	// bare domain object rather than the proto response message: it names the
	// single field to unwrap from the REST/gRPC response before unmarshalling.
	// Used by the organization / org SSO / SCIM / org domain methods, whose
	// hand-written signatures predate the proto (server 2.4.0 added the RPCs)
	// and are kept so adding the transports is not a breaking change.
	// Leave empty when the response is read whole (a message, or a paginated
	// list the domain type already mirrors).
	responseUnwrap string

	// restResponse builds the proto response message the REST body is decoded
	// into, for methods whose out is a hand-written domain type rather than the
	// proto message itself. grpc-gateway emits int64 as a JSON string, which
	// doREST only handles for proto.Message targets (via protojson).
	restResponse func() proto.Message

	// restMethod / restPath; empty restPath means rest-unsupported.
	restMethod string
	restPath   string
	restBody   interface{}

	// grpcCall invokes the admin stub; nil means grpc-unsupported.
	grpcCall func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error)
}

// adminAuthHeaders merges the admin-secret header onto the client defaults.
func (c *AuthorizerAdminClient) adminAuthHeaders() map[string]string {
	h := mergeHeaders(c.ExtraHeaders, map[string]string{adminSecretHeader: c.AdminSecret})
	return h
}

// supported reports which protocols a spec implements, for error messages.
func (s adminMethodSpec) supported() string {
	var got []string
	if s.graphql != nil {
		got = append(got, "graphql")
	}
	if s.restPath != "" {
		got = append(got, "rest")
	}
	if s.grpcCall != nil {
		got = append(got, "grpc")
	}
	return strings.Join(got, " or ")
}

// execute dispatches an admin method over the client's selected protocol and
// unmarshals the result into out (a pointer). Calling a method on a protocol it
// does not support returns a clear error before any network call.
func (c *AuthorizerAdminClient) execute(spec adminMethodSpec, out interface{}) error {
	switch c.Protocol {
	case ProtocolREST:
		if spec.restPath == "" {
			return unsupportedProtocol(spec.name, c.Protocol, spec.supported())
		}
		if spec.restResponse != nil {
			msg := spec.restResponse()
			if err := doREST(c.AuthorizerURL, spec.restMethod, spec.restPath, spec.restBody, c.ExtraHeaders, map[string]string{adminSecretHeader: c.AdminSecret}, msg); err != nil {
				return err
			}
			return unwrapProto(msg, spec.responseUnwrap, out)
		}
		return doREST(c.AuthorizerURL, spec.restMethod, spec.restPath, spec.restBody, c.ExtraHeaders, map[string]string{adminSecretHeader: c.AdminSecret}, out)

	case ProtocolGRPC:
		if spec.grpcCall == nil {
			return unsupportedProtocol(spec.name, c.Protocol, spec.supported())
		}
		conn, err := grpcDial(c.AuthorizerURL, c.GRPCEndpoint)
		if err != nil {
			return err
		}
		defer conn.Close()

		ctx := outgoingContext(context.Background(), c.adminAuthHeaders())
		cli := authorizerv1.NewAuthorizerAdminServiceClient(conn)
		resp, err := spec.grpcCall(ctx, cli)
		if err != nil {
			return err
		}
		return unwrapProto(resp, spec.responseUnwrap, out)

	default: // ProtocolGraphQL
		if spec.graphql == nil {
			return unsupportedProtocol(spec.name, c.Protocol, spec.supported())
		}
		bytesData, err := c.executeGraphQL(spec.graphql)
		if err != nil {
			return err
		}
		if out == nil {
			return nil
		}
		var res map[string]json.RawMessage
		if err := json.Unmarshal(bytesData, &res); err != nil {
			return err
		}
		field, ok := res[spec.graphqlField]
		if !ok {
			return nil
		}
		if spec.graphqlWrap != "" {
			wrapped, err := json.Marshal(map[string]json.RawMessage{spec.graphqlWrap: field})
			if err != nil {
				return err
			}
			field = wrapped
		}
		return json.Unmarshal(field, out)
	}
}

// executeGraphQL runs an admin GraphQL request, attaching the admin-secret
// header and the CSRF Origin header (reusing the user client's transport).
func (c *AuthorizerAdminClient) executeGraphQL(req *GraphQLRequest) ([]byte, error) {
	uc := &AuthorizerClient{
		AuthorizerURL: c.AuthorizerURL,
		ExtraHeaders:  c.ExtraHeaders,
	}
	return uc.ExecuteGraphQL(req, map[string]string{adminSecretHeader: c.AdminSecret})
}

// unwrapProto converts a proto response message into out, optionally pulling a
// single named field out of it first. Marshalling the proto struct with
// encoding/json (not protojson) is deliberate: it emits int64 as a JSON number,
// which the hand-written domain types can decode. A missing field leaves out at
// its zero value, matching how the graphql path treats an absent field.
func unwrapProto(msg interface{}, field string, out interface{}) error {
	if out == nil {
		return nil
	}
	if field == "" {
		return remarshal(msg, out)
	}
	var envelope map[string]json.RawMessage
	if err := remarshal(msg, &envelope); err != nil {
		return err
	}
	raw, ok := envelope[field]
	if !ok {
		return nil
	}
	return json.Unmarshal(raw, out)
}
