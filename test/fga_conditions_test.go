package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/authorizerdev/authorizer-go/v2"
	authorizerv1 "github.com/authorizerdev/authorizer-proto-go/authorizer/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// conditionModel extends the shared test model with a time-limited viewer
// grant. Requires a server newer than 2.4.9 (conditional tuples).
const conditionModel = `model
  schema 1.1
type user
type document
  relations
    define can_view: [user, user with non_expired_grant]
    define can_edit: [user]
condition non_expired_grant(current_time: timestamp, grant_time: timestamp, grant_duration: duration) {
  current_time < grant_time + grant_duration
}`

// TestFgaConditions drives conditional tuples through every protocol: an
// admin writes and reads back a time-limited grant, the server's clock decides
// expiry for the user's checks, and the server's write rules are surfaced.
func TestFgaConditions(t *testing.T) {
	c := testClient(t)
	email := uniqueEmail()
	signUp(t, c, email)
	loginRes := login(t, c, email)
	userID := loginRes.User.ID
	headers := map[string]string{"Authorization": "Bearer " + authorizer.StringValue(loginRes.AccessToken)}

	grant := func(start time.Time, d string) *authorizerv1.FgaCondition {
		ctx, err := structpb.NewStruct(map[string]interface{}{
			"grant_time": start.UTC().Format(time.RFC3339), "grant_duration": d,
		})
		if err != nil {
			t.Fatal(err)
		}
		return &authorizerv1.FgaCondition{Name: "non_expired_grant", Context: ctx}
	}

	for _, p := range adminProtocols {
		p := p
		t.Run(string(p), func(t *testing.T) {
			admin := adminClient(t, p)
			if _, err := admin.FgaWriteModel(&authorizerv1.FgaWriteModelRequest{Dsl: conditionModel}); err != nil {
				skipIfFgaUnavailable(t, err)
				t.Fatalf("FgaWriteModel: %v", err)
			}
			fresh := fmt.Sprintf("document:cond-fresh-%s", p)
			stale := fmt.Sprintf("document:cond-stale-%s", p)
			_, err := admin.FgaWriteTuples(&authorizerv1.FgaWriteTuplesRequest{Tuples: []*authorizerv1.FgaTupleInput{
				{User: "user:" + userID, Relation: "can_view", Object: fresh, Condition: grant(time.Now(), "10m")},
				{User: "user:" + userID, Relation: "can_view", Object: stale, Condition: grant(time.Now().Add(-time.Hour), "10m")},
			}})
			if err != nil {
				t.Fatalf("FgaWriteTuples with condition: %v", err)
			}

			read, err := admin.FgaReadTuples(&authorizerv1.FgaReadTuplesRequest{Object: &fresh})
			if err != nil {
				t.Fatalf("FgaReadTuples: %v", err)
			}
			if len(read.Tuples) != 1 || read.Tuples[0].GetCondition().GetName() != "non_expired_grant" ||
				read.Tuples[0].GetCondition().GetContext().AsMap()["grant_duration"] != "10m" {
				t.Fatalf("FgaReadTuples: condition not returned: %+v", read.Tuples)
			}

			res, err := protocolClient(t, p).CheckPermissions(&authorizer.CheckPermissionsRequest{
				Checks: []*authorizer.PermissionCheckInput{
					{Relation: "can_view", Object: fresh},
					{Relation: "can_view", Object: stale},
				},
			}, headers)
			if err != nil {
				t.Fatalf("CheckPermissions: %v", err)
			}
			if len(res.Results) != 2 || !res.Results[0].Allowed || res.Results[1].Allowed {
				t.Fatalf("expected fresh allowed and stale denied, got %+v", res.Results)
			}

			_, err = protocolClient(t, p).CheckPermissions(&authorizer.CheckPermissionsRequest{
				Checks:  []*authorizer.PermissionCheckInput{{Relation: "can_view", Object: stale}},
				Context: map[string]interface{}{"current_time": "2000-01-01T00:00:00Z"},
			}, headers)
			if err == nil || !strings.Contains(err.Error(), "super-admin") {
				t.Fatalf("a user's Context must be rejected, got %v", err)
			}

			missing, _ := structpb.NewStruct(map[string]interface{}{"grant_time": time.Now().UTC().Format(time.RFC3339)})
			_, err = admin.FgaWriteTuples(&authorizerv1.FgaWriteTuplesRequest{Tuples: []*authorizerv1.FgaTupleInput{
				{User: "user:" + userID, Relation: "can_view", Object: fresh + "-x",
					Condition: &authorizerv1.FgaCondition{Name: "non_expired_grant", Context: missing}},
			}})
			if err == nil || !strings.Contains(err.Error(), "missing parameters") {
				t.Fatalf("an unbound condition parameter must be rejected, got %v", err)
			}
		})
	}
}
