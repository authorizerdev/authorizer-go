package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// PermissionCheckInput is one permission to evaluate: "does the subject have
// Relation on Object?".
type PermissionCheckInput struct {
	Relation string `json:"relation"`
	Object   string `json:"object"`
	// ContextualTuples are evaluated for this check only and not persisted.
	ContextualTuples []*FgaTupleInput `json:"contextual_tuples,omitempty"`
}

// CheckPermissionsRequest evaluates one or more permission checks in a single
// call. The subject defaults to the authenticated caller.
type CheckPermissionsRequest struct {
	Checks []*PermissionCheckInput `json:"checks"`
	// User optionally overrides the subject ("type:id" or bare id → "user:id").
	// Honored only for super-admin callers or when it equals the caller's own
	// token subject; rejected otherwise.
	User string `json:"user,omitempty"`
}

// PermissionCheckResult is the outcome of one permission check, echoing the
// checked Relation/Object pair so batch results are self-describing (and
// positionally aligned with the request's Checks).
type PermissionCheckResult struct {
	Relation string `json:"relation"`
	Object   string `json:"object"`
	Allowed  bool   `json:"allowed"`
}

// CheckPermissionsResponse carries one result per supplied check, in order.
type CheckPermissionsResponse struct {
	Results []*PermissionCheckResult `json:"results"`
}

// CheckPermissions performs the check_permissions query, evaluating one or
// more relation/object checks for the authenticated caller in a single round
// trip. Results come back in the same order as the supplied Checks, each
// echoing its relation/object pair. headers must carry the caller's bearer
// token or session cookie so the server can pin the subject.
func (c *AuthorizerClient) CheckPermissions(req *CheckPermissionsRequest, headers map[string]string) (*CheckPermissionsResponse, error) {
	var res CheckPermissionsResponse
	err := c.execute(methodSpec{
		name: "CheckPermissions",
		graphql: &GraphQLRequest{
			Query:     `query checkPermissions($data: CheckPermissionsInput!){check_permissions(params: $data) { results { relation object allowed } } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "check_permissions",
		restMethod:   http.MethodPost,
		restPath:     "/v1/check_permissions",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.CheckPermissionsResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.CheckPermissionsRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.CheckPermissions(ctx, &in)
		},
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
