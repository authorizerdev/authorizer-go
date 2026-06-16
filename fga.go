package authorizer

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
//
// The individual operations live in check_permissions.go and
// list_permissions.go.

// FgaTupleInput is a single relationship tuple: User is related to Object via
// Relation. Used to pass contextual tuples that are evaluated for one check
// only and never persisted.
type FgaTupleInput struct {
	User     string `json:"user"`
	Relation string `json:"relation"`
	Object   string `json:"object"`
}
