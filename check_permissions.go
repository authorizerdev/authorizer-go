package authorizer

import "encoding/json"

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
