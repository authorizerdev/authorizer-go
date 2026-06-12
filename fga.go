package authorizer

import "encoding/json"

// This file implements the client-facing fine-grained authorization (FGA)
// surface backed by Authorizer's embedded OpenFGA engine. Only the read-side
// operations a relying party needs are exposed here — checking permissions and
// listing accessible objects. Authoring the model and relationship tuples is an
// admin concern handled from the dashboard / `_fga_*` admin GraphQL API, and is
// intentionally not part of this SDK.
//
// For every operation the subject (the "user" being checked) defaults to the
// authenticated caller and is pinned server-side from the request headers
// (bearer token or session cookie) — so headers are required. The optional
// `User` override ("type:id", or a bare id treated as "user:<id>") is honored
// only when the caller is a super-admin OR it equals the caller's own token
// subject; anything else is rejected by the server — never silently ignored.

// FgaTupleInput is a single relationship tuple: User is related to Object via
// Relation. Used to pass contextual tuples that are evaluated for one check
// only and never persisted.
type FgaTupleInput struct {
	User     string `json:"user"`
	Relation string `json:"relation"`
	Object   string `json:"object"`
}

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

// ListPermissionsRequest enumerates the objects of ObjectType the subject
// holds Relation on. Subject resolution follows the same rules as
// CheckPermissionsRequest.User.
type ListPermissionsRequest struct {
	Relation   string `json:"relation"`
	ObjectType string `json:"object_type"`
	// User optionally overrides the subject; same trust rules as
	// CheckPermissionsRequest.User.
	User string `json:"user,omitempty"`
}

// ListPermissionsResponse lists the fully-qualified object ids (e.g.
// "document:1") the subject holds the queried permission on.
type ListPermissionsResponse struct {
	Objects []string `json:"objects"`
}

// CheckPermissions performs the check_permissions query, evaluating one or
// more relation/object checks for the authenticated caller in a single round
// trip. Results come back in the same order as the supplied Checks, each
// echoing its relation/object pair. headers must carry the caller's bearer
// token or session cookie so the server can pin the subject.
func (c *AuthorizerClient) CheckPermissions(req *CheckPermissionsRequest, headers map[string]string) (*CheckPermissionsResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query:     `query checkPermissions($data: CheckPermissionsInput!){check_permissions(params: $data) { results { relation object allowed } } }`,
		Variables: map[string]interface{}{"data": req},
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]*CheckPermissionsResponse
	if err := json.Unmarshal(bytesData, &res); err != nil {
		return nil, err
	}
	return res["check_permissions"], nil
}

// ListPermissions performs the list_permissions query, returning the
// fully-qualified ids of objects of ObjectType the authenticated caller holds
// Relation on. headers must carry the caller's credentials.
func (c *AuthorizerClient) ListPermissions(req *ListPermissionsRequest, headers map[string]string) (*ListPermissionsResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query:     `query listPermissions($data: ListPermissionsInput!){list_permissions(params: $data) { objects } }`,
		Variables: map[string]interface{}{"data": req},
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]*ListPermissionsResponse
	if err := json.Unmarshal(bytesData, &res); err != nil {
		return nil, err
	}
	return res["list_permissions"], nil
}
