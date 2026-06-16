// Manual end-to-end smoke test for the Authorizer Go SDK.
//
// Exercises the public client (meta/signup/login/profile) and the admin client
// (users/webhooks/FGA) over the protocol you pick.
//
// Run against a local server (defaults shown):
//
//	AUTHORIZER_URL=http://localhost:8080 \
//	CLIENT_ID=test-client \
//	ADMIN_SECRET=admin \
//	PROTOCOL=graphql \   # graphql | rest | grpc
//	go run ./examples/manual
//
// gRPC listens on its own port (default :9091). The SDK derives it from the
// HTTP host; override with GRPC_ENDPOINT=host:port if needed. For plaintext
// gRPC the server must run with --grpc-insecure=true.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/authorizerdev/authorizer-go"
	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
)

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func protocol() authorizer.Protocol {
	switch env("PROTOCOL", "graphql") {
	case "rest":
		return authorizer.ProtocolREST
	case "grpc":
		return authorizer.ProtocolGRPC
	default:
		return authorizer.ProtocolGraphQL
	}
}

// step prints a labelled result; it never aborts so the whole flow runs.
func step(label string, v any, err error) {
	if err != nil {
		fmt.Printf("✗ %-22s error: %v\n", label, err)
		return
	}
	fmt.Printf("✓ %-22s %+v\n", label, v)
}

func main() {
	url := env("AUTHORIZER_URL", "http://localhost:8080")
	clientID := env("CLIENT_ID", "test-client")
	adminSecret := env("ADMIN_SECRET", "admin")
	p := protocol()
	fmt.Printf("== Authorizer Go SDK manual test ==\nurl=%s protocol=%s\n\n", url, p)

	userOpts := []authorizer.ClientOption{authorizer.WithProtocol(p)}
	adminOpts := []authorizer.AdminClientOption{authorizer.WithAdminProtocol(p)}
	if g := os.Getenv("GRPC_ENDPOINT"); g != "" {
		userOpts = append(userOpts, authorizer.WithGRPCEndpoint(g))
		adminOpts = append(adminOpts, authorizer.WithAdminGRPCEndpoint(g))
	}

	// ---- Public client ----
	c, err := authorizer.NewAuthorizerClient(clientID, url, "", nil, userOpts...)
	if err != nil {
		fmt.Println("failed to build client:", err)
		os.Exit(1)
	}

	meta, err := c.GetMetaData()
	step("GetMetaData", meta, err)

	email := fmt.Sprintf("go-manual-%d@example.com", time.Now().UnixNano())
	su, err := c.SignUp(&authorizer.SignUpRequest{
		Email:           authorizer.NewStringRef(email),
		Password:        "Test@12345",
		ConfirmPassword: "Test@12345",
	})
	step("SignUp", su, err)

	li, err := c.Login(&authorizer.LoginRequest{
		Email:    authorizer.NewStringRef(email),
		Password: "Test@12345",
	})
	step("Login", li, err)

	if li != nil && li.AccessToken != nil {
		bearer := map[string]string{"Authorization": "Bearer " + *li.AccessToken}
		prof, err := c.GetProfile(bearer)
		step("GetProfile", prof, err)
	}

	// ---- Admin client (auth via x-authorizer-admin-secret) ----
	fmt.Println("\n-- admin --")
	admin, err := authorizer.NewAuthorizerAdminClient(url, adminSecret, adminOpts...)
	if err != nil {
		fmt.Println("failed to build admin client:", err)
		return
	}

	users, err := admin.Users(&authorizerv1.UsersRequest{})
	if err != nil {
		step("Admin.Users", nil, err)
	} else {
		step("Admin.Users", fmt.Sprintf("%d user(s)", len(users.GetUsers())), nil)
	}

	const webhookEndpoint = "https://example.com/webhook"
	wh, err := admin.AddWebhook(&authorizerv1.AddWebhookRequest{
		EventName: "user.login",
		Endpoint:  webhookEndpoint,
		Enabled:   true,
	})
	step("Admin.AddWebhook", wh, err)

	whs, err := admin.Webhooks(&authorizerv1.WebhooksRequest{})
	if err != nil {
		step("Admin.Webhooks", nil, err)
	} else {
		step("Admin.Webhooks", fmt.Sprintf("%d webhook(s)", len(whs.GetWebhooks())), nil)
		// Clean up by endpoint: the server appends a "-<timestamp>" suffix to
		// event_name (so it is not a stable key), but endpoint is stored verbatim.
		for _, w := range whs.GetWebhooks() {
			if w.GetEndpoint() == webhookEndpoint {
				_, derr := admin.DeleteWebhook(&authorizerv1.DeleteWebhookRequest{Id: w.GetId()})
				step("Admin.DeleteWebhook", w.GetId(), derr)
			}
		}
	}

	// ---- FGA admin ----
	fmt.Println("\n-- fga admin --")
	const model = `model
  schema 1.1
type user
type document
  relations
    define viewer: [user]`
	fm, err := admin.FgaWriteModel(&authorizerv1.FgaWriteModelRequest{Dsl: model})
	step("Admin.FgaWriteModel", fm, err)

	object := fmt.Sprintf("document:%d", time.Now().UnixNano()) // unique so re-runs don't collide
	wt, err := admin.FgaWriteTuples(&authorizerv1.FgaWriteTuplesRequest{
		Tuples: []*authorizerv1.FgaTupleInput{
			{User: "user:alice", Relation: "viewer", Object: object},
		},
	})
	step("Admin.FgaWriteTuples", wt, err)

	rt, err := admin.FgaReadTuples(&authorizerv1.FgaReadTuplesRequest{})
	if err != nil {
		step("Admin.FgaReadTuples", nil, err)
	} else {
		step("Admin.FgaReadTuples", fmt.Sprintf("%d tuple(s)", len(rt.GetTuples())), nil)
	}

	fmt.Println("\ndone.")
}
