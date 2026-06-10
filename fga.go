package authorizer

import "encoding/json"

// This file implements the client-facing fine-grained authorization (FGA)
// surface backed by Authorizer's embedded OpenFGA engine. Only the read-side
// operations a relying party needs are exposed here — checking access and
// listing accessible objects. Authoring the model and relationship tuples is an
// admin concern handled from the dashboard / `_fga_*` admin GraphQL API, and is
// intentionally not part of this SDK.
//
// For every operation the subject (the "user" being checked) defaults to the
// authenticated caller and is pinned server-side from the request headers
// (bearer token or session cookie) — so headers are required. The optional
// `User` override is honored only for super-admin callers and rejected
// otherwise.

// FgaTupleInput is a single relationship tuple: User is related to Object via
// Relation. Used to pass contextual tuples that are evaluated for one check
// only and never persisted.
type FgaTupleInput struct {
	User     string `json:"user"`
	Relation string `json:"relation"`
	Object   string `json:"object"`
}

// FgaCheckRequest asks "is the caller related to Object via Relation?".
type FgaCheckRequest struct {
	Relation string `json:"relation"`
	Object   string `json:"object"`
	// ContextualTuples are evaluated for this check only and not persisted.
	ContextualTuples []*FgaTupleInput `json:"contextual_tuples,omitempty"`
	// User optionally overrides the subject ("type:id" or bare id → "user:id").
	// Honored only for super-admin callers; rejected otherwise.
	User string `json:"user,omitempty"`
}

// FgaCheckResponse is the result of a single relationship check.
type FgaCheckResponse struct {
	Allowed bool `json:"allowed"`
}

// FgaCheckPair is one relation/object pair within a batch check.
type FgaCheckPair struct {
	Relation         string           `json:"relation"`
	Object           string           `json:"object"`
	ContextualTuples []*FgaTupleInput `json:"contextual_tuples,omitempty"`
	User             string           `json:"user,omitempty"`
}

// FgaBatchCheckRequest evaluates multiple relation/object pairs in one call for
// the authenticated caller.
type FgaBatchCheckRequest struct {
	Checks []*FgaCheckPair `json:"checks"`
}

// FgaBatchCheckResponse holds the results of a batch check, positionally
// aligned with the Checks supplied in the request.
type FgaBatchCheckResponse struct {
	Results []*FgaCheckResponse `json:"results"`
}

// FgaListObjectsRequest enumerates objects of ObjectType the caller relates to
// via Relation.
type FgaListObjectsRequest struct {
	Relation   string `json:"relation"`
	ObjectType string `json:"object_type"`
	User       string `json:"user,omitempty"`
}

// FgaListObjectsResponse lists fully-qualified object ids (e.g. "document:1")
// the caller relates to.
type FgaListObjectsResponse struct {
	Objects []string `json:"objects"`
}

// FgaCheck performs the fga_check query and returns whether the authenticated
// caller has Relation on Object. headers must carry the caller's bearer token
// or session cookie so the server can pin the principal.
func (c *AuthorizerClient) FgaCheck(req *FgaCheckRequest, headers map[string]string) (*FgaCheckResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query:     `query fgaCheck($data: FgaCheckInput!){fga_check(params: $data) { allowed } }`,
		Variables: map[string]interface{}{"data": req},
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]*FgaCheckResponse
	if err := json.Unmarshal(bytesData, &res); err != nil {
		return nil, err
	}
	return res["fga_check"], nil
}

// FgaBatchCheck performs the fga_batch_check query, returning one result per
// pair in Checks (same order). headers must carry the caller's credentials.
func (c *AuthorizerClient) FgaBatchCheck(req *FgaBatchCheckRequest, headers map[string]string) (*FgaBatchCheckResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query:     `query fgaBatchCheck($data: FgaBatchCheckInput!){fga_batch_check(params: $data) { results { allowed } } }`,
		Variables: map[string]interface{}{"data": req},
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]*FgaBatchCheckResponse
	if err := json.Unmarshal(bytesData, &res); err != nil {
		return nil, err
	}
	return res["fga_batch_check"], nil
}

// FgaListObjects performs the fga_list_objects query, returning the
// fully-qualified ids of objects of ObjectType the caller relates to via
// Relation. headers must carry the caller's credentials.
func (c *AuthorizerClient) FgaListObjects(req *FgaListObjectsRequest, headers map[string]string) (*FgaListObjectsResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query:     `query fgaListObjects($data: FgaListObjectsInput!){fga_list_objects(params: $data) { objects } }`,
		Variables: map[string]interface{}{"data": req},
	}, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]*FgaListObjectsResponse
	if err := json.Unmarshal(bytesData, &res); err != nil {
		return nil, err
	}
	return res["fga_list_objects"], nil
}
