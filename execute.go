package authorizer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// methodSpec declares how one public client method is carried over each
// protocol. A protocol whose field is empty/nil is unsupported and yields a
// clear "<name> not available over <protocol>; use <supported>" error before
// any network call — mirroring the admin client's availability map. As of
// 2.3.0 every public RPC is implemented over graphql + rest + grpc.
//
// For grpc/rest the spec returns proto response messages. The response envelope
// is now flat (no per-RPC wrapper), so the bare proto domain message maps
// directly onto the SDK's response struct. Because proto JSON (protojson +
// grpc-gateway REST) emits int64 fields as strings, REST bodies are decoded
// with protojson before being round-tripped onto the SDK type.
type methodSpec struct {
	// name is the human method name for error messages, e.g. "Login".
	name string

	// graphql is the prebuilt GraphQL request; nil means gql-unsupported.
	graphql *GraphQLRequest
	// graphqlField is the top-level response field to unwrap, e.g. "login".
	graphqlField string

	// restMethod / restPath are the typed REST endpoint, e.g.
	// (http.MethodPost, "/v1/login"). An empty restPath means rest-unsupported.
	restMethod string
	restPath   string
	// restBody is the JSON payload for POST endpoints (nil for GET).
	restBody interface{}
	// restResp returns a fresh proto response message to protojson-unmarshal the
	// REST body into (so int64 string fields decode correctly).
	restResp func() proto.Message

	// grpcCall invokes the generated stub and returns the proto response; nil
	// means grpc-unsupported.
	grpcCall func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error)
}

// supported reports which protocols a spec implements, for error messages.
func (s methodSpec) supported() string {
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

// projectProto round-trips a flat proto response message onto the SDK's
// response struct (out).
func (s methodSpec) projectProto(resp interface{}, out interface{}) error {
	return remarshal(resp, out)
}

// execute dispatches a public client method over the client's selected
// protocol and unmarshals the result into out (a pointer). Calling a method on
// a protocol it does not support returns a clear error before any network call.
func (c *AuthorizerClient) execute(spec methodSpec, headers map[string]string, out interface{}) error {
	switch c.Protocol {
	case ProtocolREST:
		if spec.restPath == "" {
			return unsupportedProtocol(spec.name, c.Protocol, spec.supported())
		}
		// Decode proto-typed REST responses with protojson (int64 serialize as
		// strings over REST), then map the flat proto message onto the SDK type.
		if out == nil || spec.restResp == nil {
			return c.executeREST(spec.restMethod, spec.restPath, spec.restBody, headers, nil)
		}
		resp := spec.restResp()
		if err := c.executeREST(spec.restMethod, spec.restPath, spec.restBody, headers, resp); err != nil {
			return err
		}
		return spec.projectProto(resp, out)

	case ProtocolGRPC:
		if spec.grpcCall == nil {
			return unsupportedProtocol(spec.name, c.Protocol, spec.supported())
		}
		conn, err := grpcDial(c.AuthorizerURL, c.GRPCEndpoint)
		if err != nil {
			return err
		}
		defer conn.Close()

		ctx := outgoingContext(context.Background(), mergeHeaders(c.ExtraHeaders, headers))
		cli := authorizerv1.NewAuthorizerServiceClient(conn)
		resp, err := spec.grpcCall(ctx, cli)
		if err != nil {
			return err
		}
		if out == nil {
			return nil
		}
		return spec.projectProto(resp, out)

	default: // ProtocolGraphQL
		if spec.graphql == nil {
			return unsupportedProtocol(spec.name, c.Protocol, spec.supported())
		}
		bytesData, err := c.ExecuteGraphQL(spec.graphql, headers)
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
		return json.Unmarshal(field, out)
	}
}

// remarshal JSON round-trips a proto response message into the SDK's own
// response struct, so the public method signatures stay identical across
// protocols. Proto messages carry json tags matching the GraphQL/REST shapes.
func remarshal(from, to interface{}) error {
	if to == nil {
		return nil
	}
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, to)
}

// mergeHeaders merges per-call headers over the client's default headers.
func mergeHeaders(base, override map[string]string) map[string]string {
	if len(base) == 0 {
		return override
	}
	merged := make(map[string]string, len(base)+len(override))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range override {
		merged[k] = v
	}
	return merged
}

// unsupportedProtocol builds the standard "<method> not available over
// <protocol>" error used by methods that only exist on a subset of protocols.
func unsupportedProtocol(method string, p Protocol, supported string) error {
	return fmt.Errorf("%s not available over %s; use %s", method, p, supported)
}

// grpcOpts is a convenience alias so spec closures stay short.
type grpcOpts = []grpc.CallOption
