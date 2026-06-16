package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
	"google.golang.org/protobuf/proto"
)

// ListPermissionsRequest enumerates what the subject can access. With both
// Relation and ObjectType set it answers "which <ObjectType>s can I
// <Relation>?". Either or both filters may be omitted: every matching
// (type, relation) pair of the active model is then enumerated, so an empty
// request returns ALL permissions the subject holds. Subject resolution
// follows the same rules as CheckPermissionsRequest.User.
type ListPermissionsRequest struct {
	Relation   string `json:"relation,omitempty"`
	ObjectType string `json:"object_type,omitempty"`
	// User optionally overrides the subject; same trust rules as
	// CheckPermissionsRequest.User.
	User string `json:"user,omitempty"`
}

// Permission is one (object, relation) pair the subject holds.
type Permission struct {
	Object   string `json:"object"`
	Relation string `json:"relation"`
}

// ListPermissionsResponse lists what the subject can access. Objects is the
// distinct fully-qualified object ids (e.g. "document:1"); Permissions carries
// the (object, relation) detail — relevant when no Relation filter was
// supplied. Truncated is true when the result was capped (1000 entries) and
// more permissions exist.
type ListPermissionsResponse struct {
	Objects     []string      `json:"objects"`
	Permissions []*Permission `json:"permissions"`
	Truncated   bool          `json:"truncated"`
}

// ListPermissions performs the list_permissions query, returning the
// fully-qualified ids of objects of ObjectType the authenticated caller holds
// Relation on (or, with filters omitted, everything the caller holds).
// headers must carry the caller's credentials.
func (c *AuthorizerClient) ListPermissions(req *ListPermissionsRequest, headers map[string]string) (*ListPermissionsResponse, error) {
	var res ListPermissionsResponse
	err := c.execute(methodSpec{
		name: "ListPermissions",
		graphql: &GraphQLRequest{
			Query:     `query listPermissions($data: ListPermissionsInput!){list_permissions(params: $data) { objects permissions { object relation } truncated } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "list_permissions",
		restMethod:   http.MethodPost,
		restPath:     "/v1/list_permissions",
		restBody:     req,
		restResp:     func() proto.Message { return &authorizerv1.ListPermissionsResponse{} },
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerServiceClient) (interface{}, error) {
			var in authorizerv1.ListPermissionsRequest
			if err := remarshal(req, &in); err != nil {
				return nil, err
			}
			return cli.ListPermissions(ctx, &in)
		},
	}, headers, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
