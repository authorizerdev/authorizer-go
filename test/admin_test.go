package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/authorizerdev/authorizer-go"
	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
)

const adminSecret = "admin"

// adminProtocols are the transports the admin surface supports end-to-end.
var adminProtocols = []authorizer.Protocol{
	authorizer.ProtocolGraphQL,
	authorizer.ProtocolREST,
	authorizer.ProtocolGRPC,
}

// adminClient builds an admin client pinned to the given protocol.
func adminClient(t *testing.T, p authorizer.Protocol) *authorizer.AuthorizerAdminClient {
	t.Helper()
	c, err := authorizer.NewAuthorizerAdminClient(authorizerURL, adminSecret,
		authorizer.WithAdminProtocol(p),
		authorizer.WithAdminExtraHeaders(map[string]string{"Origin": authorizerURL}),
	)
	if err != nil {
		t.Fatalf("failed to create %s admin client: %v", p, err)
	}
	return c
}

func TestNewAuthorizerAdminClientValidation(t *testing.T) {
	if _, err := authorizer.NewAuthorizerAdminClient("", "secret"); err == nil {
		t.Error("expected error for missing authorizerURL")
	}
	if _, err := authorizer.NewAuthorizerAdminClient(authorizerURL, ""); err == nil {
		t.Error("expected error for missing adminSecret")
	}
	c, err := authorizer.NewAuthorizerAdminClient(authorizerURL, adminSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Protocol != authorizer.ProtocolGraphQL {
		t.Errorf("expected default admin protocol graphql, got %q", c.Protocol)
	}
}

// TestAdminUsersAcrossProtocols lists users over each supported protocol.
func TestAdminUsersAcrossProtocols(t *testing.T) {
	for _, p := range adminProtocols {
		p := p
		t.Run(string(p), func(t *testing.T) {
			c := adminClient(t, p)
			res, err := c.Users(&authorizerv1.UsersRequest{
				Pagination: &authorizerv1.PaginationRequest{Page: 1, Limit: 10},
			})
			if err != nil {
				t.Fatalf("[%s] Users failed: %v", p, err)
			}
			if res == nil {
				t.Fatalf("[%s] Users returned nil response", p)
			}
		})
	}
}

// TestAdminMetaProtocolAvailability verifies AdminMeta works over rest+grpc and
// returns a clear error over graphql (which has no _admin_meta-shaped op here).
func TestAdminMetaProtocolAvailability(t *testing.T) {
	// rest + grpc: supported
	for _, p := range []authorizer.Protocol{authorizer.ProtocolREST, authorizer.ProtocolGRPC} {
		c := adminClient(t, p)
		if _, err := c.AdminMeta(); err != nil {
			t.Fatalf("[%s] AdminMeta failed: %v", p, err)
		}
	}

	// graphql: unsupported → clear error, no network 404
	c := adminClient(t, authorizer.ProtocolGraphQL)
	_, err := c.AdminMeta()
	if err == nil {
		t.Fatal("expected AdminMeta to error over graphql")
	}
	if !strings.Contains(err.Error(), "not available over graphql") {
		t.Errorf("expected clear unsupported-protocol error, got %v", err)
	}
}

// TestAdminGqlOnlyExtras verifies the gql-only methods error over rest+grpc.
func TestAdminGqlOnlyExtras(t *testing.T) {
	for _, p := range []authorizer.Protocol{authorizer.ProtocolREST, authorizer.ProtocolGRPC} {
		c := adminClient(t, p)
		if _, err := c.GenerateJWTKeys(&authorizer.GenerateJWTKeysRequest{Type: "HS256"}); err == nil {
			t.Errorf("[%s] expected GenerateJWTKeys to error (gql-only)", p)
		} else if !strings.Contains(err.Error(), "use graphql") {
			t.Errorf("[%s] expected 'use graphql' hint, got %v", p, err)
		}
	}
}

// webhookEventByProtocol assigns each protocol a DISTINCT, server-valid webhook
// event name so the cross-protocol loop never trips the
// authorizer_webhooks.event_name UNIQUE constraint (event names are validated
// against a fixed server-side enum, so a random suffix is rejected).
var webhookEventByProtocol = map[authorizer.Protocol]string{
	authorizer.ProtocolGraphQL: "user.login",
	authorizer.ProtocolREST:    "user.signup",
	authorizer.ProtocolGRPC:    "user.created",
}

// TestAdminWebhookLifecycle exercises add → list (locate) → delete over each
// supported protocol. Each protocol uses a distinct valid event_name, and the
// created webhook is always deleted at the end (even on failure) so it does not
// pollute later runs.
func TestAdminWebhookLifecycle(t *testing.T) {
	for _, p := range adminProtocols {
		p := p
		t.Run(string(p), func(t *testing.T) {
			c := adminClient(t, p)
			eventName := webhookEventByProtocol[p]

			added, err := c.AddWebhook(&authorizerv1.AddWebhookRequest{
				EventName: eventName,
				Endpoint:  "https://example.com/webhook",
				Enabled:   true,
			})
			if err != nil {
				t.Fatalf("[%s] AddWebhook failed: %v", p, err)
			}
			if added == nil {
				t.Fatalf("[%s] AddWebhook returned nil", p)
			}

			// Locate the webhook by its event_name to get its id (the list is
			// paginated — an explicit limit is required or it returns empty),
			// then schedule deletion so the fixture is cleaned up. The storage
			// layer appends a "-<unix-ts>" suffix to the stored event_name, so
			// match by prefix rather than exact equality.
			list, err := c.Webhooks(&authorizerv1.WebhooksRequest{
				Pagination: &authorizerv1.PaginationRequest{Page: 1, Limit: 1000},
			})
			if err != nil {
				t.Fatalf("[%s] Webhooks failed: %v", p, err)
			}
			var id string
			for _, w := range list.GetWebhooks() {
				if strings.HasPrefix(w.GetEventName(), eventName+"-") {
					id = w.GetId()
					break
				}
			}
			if id == "" {
				t.Fatalf("[%s] created webhook %q not found in list", p, eventName)
			}
			t.Cleanup(func() {
				if _, err := c.DeleteWebhook(&authorizerv1.DeleteWebhookRequest{Id: id}); err != nil {
					t.Errorf("[%s] DeleteWebhook cleanup failed: %v", p, err)
				}
			})

			// GetWebhook by id (round-trips the int64 created_at/updated_at over
			// REST via protojson).
			if _, err := c.GetWebhook(&authorizerv1.GetWebhookRequest{Id: id}); err != nil {
				t.Fatalf("[%s] GetWebhook failed: %v", p, err)
			}
		})
	}
}

// TestAdminFgaModelAndTuples drives the FGA admin flow: write a model, write
// tuples, then read / list / expand. Each protocol writes a DISTINCT object
// ("document:<protocol>") so re-running the same tuple across the loop does not
// trip openfga's "tuple already exists" guard on a shared store. FgaReset
// (destructive) is NOT run here — it runs once at the very end in
// TestAdminFgaResetLast so the suite never wipes the store mid-run.
func TestAdminFgaModelAndTuples(t *testing.T) {
	const model = `model
  schema 1.1
type user
type document
  relations
    define can_view: [user]`

	for _, p := range adminProtocols {
		p := p
		t.Run(string(p), func(t *testing.T) {
			c := adminClient(t, p)
			object := fmt.Sprintf("document:%s", p)

			if _, err := c.FgaWriteModel(&authorizerv1.FgaWriteModelRequest{Dsl: model}); err != nil {
				skipIfFgaUnavailable(t, err)
				t.Fatalf("[%s] FgaWriteModel failed: %v", p, err)
			}

			// Write-once-then-verify: tolerate "already exists" so the test is
			// rerun-safe against a shared (non-ephemeral) store.
			if _, err := c.FgaWriteTuples(&authorizerv1.FgaWriteTuplesRequest{
				Tuples: []*authorizerv1.FgaTupleInput{
					{User: "user:1", Relation: "can_view", Object: object},
				},
			}); err != nil && !strings.Contains(err.Error(), "already exist") {
				t.Fatalf("[%s] FgaWriteTuples failed: %v", p, err)
			}

			if _, err := c.FgaReadTuples(&authorizerv1.FgaReadTuplesRequest{}); err != nil {
				t.Fatalf("[%s] FgaReadTuples failed: %v", p, err)
			}

			if _, err := c.FgaListUsers(&authorizerv1.FgaListUsersRequest{
				Object:   object,
				Relation: "can_view",
				UserType: "user",
			}); err != nil {
				t.Fatalf("[%s] FgaListUsers failed: %v", p, err)
			}

			if _, err := c.FgaExpand(&authorizerv1.FgaExpandRequest{
				Relation: "can_view",
				Object:   object,
			}); err != nil {
				t.Fatalf("[%s] FgaExpand failed: %v", p, err)
			}
		})
	}
}

// TestAdminFgaResetLast exercises the destructive FgaReset exactly once (over
// rest — FgaReset is a rest/grpc-only RPC, no graphql op) so it does not run
// inside the per-protocol loop and wipe the store mid-suite. The server refuses
// to reset the model while tuples exist, so every existing tuple is read and
// deleted first (covers tuples written by TestAdminFgaModelAndTuples plus any
// left over from a prior run, keeping the suite rerun-safe).
//
// DESTRUCTIVE: FgaReset clears all FGA tuples and the model on the server.
func TestAdminFgaResetLast(t *testing.T) {
	c := adminClient(t, authorizer.ProtocolREST)

	read, err := c.FgaReadTuples(&authorizerv1.FgaReadTuplesRequest{})
	if err != nil {
		skipIfFgaUnavailable(t, err)
		t.Fatalf("FgaReadTuples failed: %v", err)
	}
	var tuples []*authorizerv1.FgaTupleInput
	for _, tp := range read.GetTuples() {
		tuples = append(tuples, &authorizerv1.FgaTupleInput{
			User: tp.GetUser(), Relation: tp.GetRelation(), Object: tp.GetObject(),
		})
	}
	if len(tuples) > 0 {
		if _, err := c.FgaDeleteTuples(&authorizerv1.FgaDeleteTuplesRequest{Tuples: tuples}); err != nil {
			t.Fatalf("FgaDeleteTuples failed: %v", err)
		}
	}

	if _, err := c.FgaReset(); err != nil {
		t.Fatalf("FgaReset failed: %v", err)
	}
}
