package authorizer

import (
	"context"
	"net/http"

	"google.golang.org/protobuf/proto"

	authorizerv1 "github.com/authorizerdev/authorizer-proto-go/authorizer/v1"
)

// GraphQL selection-set fragments for the admin response types. Kept here so
// each method's query stays readable. Field names mirror the GraphQL schema.
const (
	adminUserFields       = UserFragment
	adminPaginationFields = `pagination { limit page offset total }`
	adminWebhookFields    = `id event_name event_description endpoint enabled headers created_at updated_at`
	adminWebhookLogFields = `id http_status response request webhook_id created_at updated_at`
	adminEmailTplFields   = `id event_name template design subject created_at updated_at`
	adminAuditLogFields   = `id actor_id actor_type actor_email action resource_type resource_id ip_address user_agent metadata created_at`
	adminVerifReqFields   = `id identifier token email expires created_at updated_at nonce redirect_uri`
	adminFgaModelFields   = `id dsl`
	adminFgaTupleFields   = `user relation object`

	adminClientFields        = `id client_id name description allowed_scopes is_active created_at updated_at`
	adminTrustedIssuerFields = `id service_account_id name issuer_url key_source_type jwks_url expected_aud subject_claim allowed_subjects issuer_type is_active spiffe_refresh_hint_seconds created_at updated_at`
	adminOrgFields           = `id name display_name enabled created_at updated_at`
	adminOrgMemberFields     = `id org_id user_id email given_name family_name roles created_at updated_at`
	adminOrgOIDCConnFields   = `id org_id name issuer_url sso_client_id scopes redirect_uri is_active created_at updated_at`
	adminOrgSAMLConnFields   = `id org_id name idp_entity_id idp_sso_url sp_entity_id acs_url attribute_mapping allow_idp_initiated is_active created_at updated_at`
	adminScimEndpointFields  = `id org_id enabled created_at updated_at`
	adminSamlSPFields        = `id org_id name entity_id acs_url sp_cert_pem name_id_format mapped_attributes allow_idp_initiated is_active created_at updated_at`
	adminSamlIdpKeyFields    = `id org_id cert_pem algorithm status created_at updated_at`
)

// ---------------------------------------------------------------------------
// 1. AdminLogin — establishes an admin session. grpc, rest, gql.
// ---------------------------------------------------------------------------

// AdminLogin validates the admin secret and establishes an admin session.
func (c *AuthorizerAdminClient) AdminLogin(req *authorizerv1.AdminLoginRequest) (*Response, error) {
	var res Response
	err := c.execute(adminMethodSpec{
		name: "AdminLogin",
		graphql: &GraphQLRequest{
			Query:     `mutation adminLogin($data: AdminLoginRequest!) { _admin_login(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_admin_login",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/login",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.AdminLogin(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 2. AdminLogout — grpc, rest (no gql).
// ---------------------------------------------------------------------------

// AdminLogout clears the admin session.
func (c *AuthorizerAdminClient) AdminLogout() (*authorizerv1.AdminLogoutResponse, error) {
	var res authorizerv1.AdminLogoutResponse
	err := c.execute(adminMethodSpec{
		name: "AdminLogout",
		graphql: &GraphQLRequest{
			Query: `mutation adminLogout { _admin_logout { message } }`,
		},
		graphqlField: "_admin_logout",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/logout",
		restBody:     &authorizerv1.AdminLogoutRequest{},
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.AdminLogout(ctx, &authorizerv1.AdminLogoutRequest{})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 3. AdminSession — grpc, rest (no gql).
// ---------------------------------------------------------------------------

// AdminSession refreshes the admin session.
func (c *AuthorizerAdminClient) AdminSession() (*authorizerv1.AdminSessionResponse, error) {
	var res authorizerv1.AdminSessionResponse
	err := c.execute(adminMethodSpec{
		name: "AdminSession",
		graphql: &GraphQLRequest{
			Query: `query adminSession { _admin_session { message } }`,
		},
		graphqlField: "_admin_session",
		restMethod:   http.MethodGet,
		restPath:     "/v1/admin/session",
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.AdminSession(ctx, &authorizerv1.AdminSessionRequest{})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 4. AdminMeta — grpc, rest (no gql).
// ---------------------------------------------------------------------------

// AdminMeta returns admin-only configuration metadata (roles, default roles,
// protected roles).
func (c *AuthorizerAdminClient) AdminMeta() (*authorizerv1.AdminMetaResponse, error) {
	var res authorizerv1.AdminMetaResponse
	err := c.execute(adminMethodSpec{
		name: "AdminMeta",
		graphql: &GraphQLRequest{
			Query: `query adminMeta { _admin_meta { roles default_roles protected_roles is_multi_factor_auth_service_enabled } }`,
		},
		graphqlField: "_admin_meta",
		graphqlWrap:  "admin_meta",
		restMethod:   http.MethodGet,
		restPath:     "/v1/admin/meta",
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.AdminMeta(ctx, &authorizerv1.AdminMetaRequest{})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 5. Users — grpc, rest, gql.
// ---------------------------------------------------------------------------

// Users returns a paginated list of users.
func (c *AuthorizerAdminClient) Users(req *authorizerv1.UsersRequest) (*authorizerv1.UsersResponse, error) {
	var res authorizerv1.UsersResponse
	err := c.execute(adminMethodSpec{
		name: "Users",
		graphql: &GraphQLRequest{
			Query:     "query users($data: ListUsersRequest) { _users(params: $data) { " + adminPaginationFields + " users { " + adminUserFields + " } } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_users",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/users",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.Users(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 6. User — grpc, rest, gql.
// ---------------------------------------------------------------------------

// User returns a single user by id or email.
func (c *AuthorizerAdminClient) User(req *authorizerv1.UserRequest) (*authorizerv1.UserResponse, error) {
	var res authorizerv1.UserResponse
	err := c.execute(adminMethodSpec{
		name: "User",
		graphql: &GraphQLRequest{
			Query:     "query user($data: GetUserRequest!) { _user(params: $data) { " + adminUserFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_user",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/user",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.User(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 7. UpdateUser — grpc, rest, gql.
// ---------------------------------------------------------------------------

// UpdateUser updates an existing user.
func (c *AuthorizerAdminClient) UpdateUser(req *authorizerv1.UpdateUserRequest) (*authorizerv1.UpdateUserResponse, error) {
	var res authorizerv1.UpdateUserResponse
	err := c.execute(adminMethodSpec{
		name: "UpdateUser",
		graphql: &GraphQLRequest{
			Query:     "mutation updateUser($data: UpdateUserRequest!) { _update_user(params: $data) { " + adminUserFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_update_user",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/update_user",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.UpdateUser(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 8. DeleteUser — grpc, rest, gql.
// ---------------------------------------------------------------------------

// DeleteUser deletes a user. DESTRUCTIVE: permanently removes the user account.
func (c *AuthorizerAdminClient) DeleteUser(req *authorizerv1.DeleteUserRequest) (*authorizerv1.DeleteUserResponse, error) {
	var res authorizerv1.DeleteUserResponse
	err := c.execute(adminMethodSpec{
		name: "DeleteUser",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteUser($data: DeleteUserRequest!) { _delete_user(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_user",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_user",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteUser(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 9. VerificationRequests — grpc, rest, gql.
// ---------------------------------------------------------------------------

// VerificationRequests returns a paginated list of pending verification requests.
func (c *AuthorizerAdminClient) VerificationRequests(req *authorizerv1.VerificationRequestsRequest) (*authorizerv1.VerificationRequestsResponse, error) {
	var res authorizerv1.VerificationRequestsResponse
	err := c.execute(adminMethodSpec{
		name: "VerificationRequests",
		graphql: &GraphQLRequest{
			Query:     "query verificationRequests($data: PaginationRequest) { _verification_requests(params: $data) { " + adminPaginationFields + " verification_requests { " + adminVerifReqFields + " } } }",
			Variables: map[string]interface{}{"data": req.GetPagination()},
		},
		graphqlField: "_verification_requests",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/verification_requests",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.VerificationRequests(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 10. RevokeAccess — grpc, rest, gql.
// ---------------------------------------------------------------------------

// RevokeAccess revokes a user's access to the application.
func (c *AuthorizerAdminClient) RevokeAccess(req *authorizerv1.RevokeAccessRequest) (*authorizerv1.RevokeAccessResponse, error) {
	var res authorizerv1.RevokeAccessResponse
	err := c.execute(adminMethodSpec{
		name: "RevokeAccess",
		graphql: &GraphQLRequest{
			Query:     `mutation revokeAccess($data: UpdateAccessRequest!) { _revoke_access(param: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_revoke_access",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/revoke_access",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.RevokeAccess(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 11. EnableAccess — grpc, rest, gql.
// ---------------------------------------------------------------------------

// EnableAccess re-enables a previously revoked user's access.
func (c *AuthorizerAdminClient) EnableAccess(req *authorizerv1.EnableAccessRequest) (*authorizerv1.EnableAccessResponse, error) {
	var res authorizerv1.EnableAccessResponse
	err := c.execute(adminMethodSpec{
		name: "EnableAccess",
		graphql: &GraphQLRequest{
			Query:     `mutation enableAccess($data: UpdateAccessRequest!) { _enable_access(param: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_enable_access",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/enable_access",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.EnableAccess(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 12. InviteMembers — grpc, rest, gql.
// ---------------------------------------------------------------------------

// InviteMembers sends invitations to one or more email addresses.
func (c *AuthorizerAdminClient) InviteMembers(req *authorizerv1.InviteMembersRequest) (*authorizerv1.InviteMembersResponse, error) {
	var res authorizerv1.InviteMembersResponse
	err := c.execute(adminMethodSpec{
		name: "InviteMembers",
		graphql: &GraphQLRequest{
			Query:     "mutation inviteMembers($data: InviteMemberRequest!) { _invite_members(params: $data) { message Users { " + adminUserFields + " } } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_invite_members",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/invite_members",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.InviteMembers(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 13. AddWebhook — grpc, rest, gql.
// ---------------------------------------------------------------------------

// AddWebhook registers a new webhook endpoint.
func (c *AuthorizerAdminClient) AddWebhook(req *authorizerv1.AddWebhookRequest) (*authorizerv1.AddWebhookResponse, error) {
	var res authorizerv1.AddWebhookResponse
	err := c.execute(adminMethodSpec{
		name: "AddWebhook",
		graphql: &GraphQLRequest{
			Query:     `mutation addWebhook($data: AddWebhookRequest!) { _add_webhook(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_add_webhook",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/add_webhook",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.AddWebhook(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 14. UpdateWebhook — grpc, rest, gql.
// ---------------------------------------------------------------------------

// UpdateWebhook updates an existing webhook.
func (c *AuthorizerAdminClient) UpdateWebhook(req *authorizerv1.UpdateWebhookRequest) (*authorizerv1.UpdateWebhookResponse, error) {
	var res authorizerv1.UpdateWebhookResponse
	err := c.execute(adminMethodSpec{
		name: "UpdateWebhook",
		graphql: &GraphQLRequest{
			Query:     `mutation updateWebhook($data: UpdateWebhookRequest!) { _update_webhook(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_update_webhook",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/update_webhook",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.UpdateWebhook(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 15. DeleteWebhook — grpc, rest, gql.
// ---------------------------------------------------------------------------

// DeleteWebhook deletes a webhook. DESTRUCTIVE: permanently removes the webhook
// and its logs.
func (c *AuthorizerAdminClient) DeleteWebhook(req *authorizerv1.DeleteWebhookRequest) (*authorizerv1.DeleteWebhookResponse, error) {
	var res authorizerv1.DeleteWebhookResponse
	err := c.execute(adminMethodSpec{
		name: "DeleteWebhook",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteWebhook($data: WebhookRequest!) { _delete_webhook(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_webhook",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_webhook",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteWebhook(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 16. GetWebhook — grpc, rest, gql.
// ---------------------------------------------------------------------------

// GetWebhook returns a single webhook by id.
func (c *AuthorizerAdminClient) GetWebhook(req *authorizerv1.GetWebhookRequest) (*authorizerv1.GetWebhookResponse, error) {
	var res authorizerv1.GetWebhookResponse
	err := c.execute(adminMethodSpec{
		name: "GetWebhook",
		graphql: &GraphQLRequest{
			Query:     "query webhook($data: WebhookRequest!) { _webhook(params: $data) { " + adminWebhookFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_webhook",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/webhook",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.GetWebhook(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 17. Webhooks — grpc, rest, gql.
// ---------------------------------------------------------------------------

// Webhooks returns a paginated list of webhooks.
func (c *AuthorizerAdminClient) Webhooks(req *authorizerv1.WebhooksRequest) (*authorizerv1.WebhooksResponse, error) {
	var res authorizerv1.WebhooksResponse
	err := c.execute(adminMethodSpec{
		name: "Webhooks",
		graphql: &GraphQLRequest{
			Query:     "query webhooks($data: PaginationRequest) { _webhooks(params: $data) { " + adminPaginationFields + " webhooks { " + adminWebhookFields + " } } }",
			Variables: map[string]interface{}{"data": req.GetPagination()},
		},
		graphqlField: "_webhooks",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/webhooks",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.Webhooks(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 18. WebhookLogs — grpc, rest, gql.
// ---------------------------------------------------------------------------

// WebhookLogs returns a paginated list of webhook delivery logs.
func (c *AuthorizerAdminClient) WebhookLogs(req *authorizerv1.WebhookLogsRequest) (*authorizerv1.WebhookLogsResponse, error) {
	var res authorizerv1.WebhookLogsResponse
	err := c.execute(adminMethodSpec{
		name: "WebhookLogs",
		graphql: &GraphQLRequest{
			Query:     "query webhookLogs($data: ListWebhookLogRequest) { _webhook_logs(params: $data) { " + adminPaginationFields + " webhook_logs { " + adminWebhookLogFields + " } } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_webhook_logs",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/webhook_logs",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.WebhookLogs(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 19. TestEndpoint — grpc, rest, gql.
// ---------------------------------------------------------------------------

// TestEndpoint sends a test event to a webhook endpoint.
func (c *AuthorizerAdminClient) TestEndpoint(req *authorizerv1.TestEndpointRequest) (*authorizerv1.TestEndpointResponse, error) {
	var res authorizerv1.TestEndpointResponse
	err := c.execute(adminMethodSpec{
		name: "TestEndpoint",
		graphql: &GraphQLRequest{
			Query:     `mutation testEndpoint($data: TestEndpointRequest!) { _test_endpoint(params: $data) { http_status response } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_test_endpoint",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/test_endpoint",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.TestEndpoint(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 20. AddEmailTemplate — grpc, rest, gql.
// ---------------------------------------------------------------------------

// AddEmailTemplate creates a new email template.
func (c *AuthorizerAdminClient) AddEmailTemplate(req *authorizerv1.AddEmailTemplateRequest) (*authorizerv1.AddEmailTemplateResponse, error) {
	var res authorizerv1.AddEmailTemplateResponse
	err := c.execute(adminMethodSpec{
		name: "AddEmailTemplate",
		graphql: &GraphQLRequest{
			Query:     `mutation addEmailTemplate($data: AddEmailTemplateRequest!) { _add_email_template(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_add_email_template",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/add_email_template",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.AddEmailTemplate(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 21. UpdateEmailTemplate — grpc, rest, gql.
// ---------------------------------------------------------------------------

// UpdateEmailTemplate updates an existing email template.
func (c *AuthorizerAdminClient) UpdateEmailTemplate(req *authorizerv1.UpdateEmailTemplateRequest) (*authorizerv1.UpdateEmailTemplateResponse, error) {
	var res authorizerv1.UpdateEmailTemplateResponse
	err := c.execute(adminMethodSpec{
		name: "UpdateEmailTemplate",
		graphql: &GraphQLRequest{
			Query:     `mutation updateEmailTemplate($data: UpdateEmailTemplateRequest!) { _update_email_template(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_update_email_template",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/update_email_template",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.UpdateEmailTemplate(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 22. DeleteEmailTemplate — grpc, rest, gql.
// ---------------------------------------------------------------------------

// DeleteEmailTemplate deletes an email template. DESTRUCTIVE: permanently
// removes the template.
func (c *AuthorizerAdminClient) DeleteEmailTemplate(req *authorizerv1.DeleteEmailTemplateRequest) (*authorizerv1.DeleteEmailTemplateResponse, error) {
	var res authorizerv1.DeleteEmailTemplateResponse
	err := c.execute(adminMethodSpec{
		name: "DeleteEmailTemplate",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteEmailTemplate($data: DeleteEmailTemplateRequest!) { _delete_email_template(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_email_template",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_email_template",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteEmailTemplate(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 23. EmailTemplates — grpc, rest, gql.
// ---------------------------------------------------------------------------

// EmailTemplates returns a paginated list of email templates.
func (c *AuthorizerAdminClient) EmailTemplates(req *authorizerv1.EmailTemplatesRequest) (*authorizerv1.EmailTemplatesResponse, error) {
	var res authorizerv1.EmailTemplatesResponse
	err := c.execute(adminMethodSpec{
		name: "EmailTemplates",
		graphql: &GraphQLRequest{
			Query:     "query emailTemplates($data: PaginationRequest) { _email_templates(params: $data) { " + adminPaginationFields + " email_templates { " + adminEmailTplFields + " } } }",
			Variables: map[string]interface{}{"data": req.GetPagination()},
		},
		graphqlField: "_email_templates",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/email_templates",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.EmailTemplates(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 24. AuditLogs — grpc, rest, gql.
// ---------------------------------------------------------------------------

// AuditLogs returns a paginated list of audit log entries.
func (c *AuthorizerAdminClient) AuditLogs(req *authorizerv1.AuditLogsRequest) (*authorizerv1.AuditLogsResponse, error) {
	var res authorizerv1.AuditLogsResponse
	err := c.execute(adminMethodSpec{
		name: "AuditLogs",
		graphql: &GraphQLRequest{
			Query:     "query auditLogs($data: ListAuditLogRequest) { _audit_logs(params: $data) { " + adminPaginationFields + " audit_logs { " + adminAuditLogFields + " } } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_audit_logs",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/audit_logs",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.AuditLogs(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 25. FgaGetModel — grpc, rest (no gql).
// ---------------------------------------------------------------------------

// FgaGetModel returns the active fine-grained authorization model as DSL.
func (c *AuthorizerAdminClient) FgaGetModel() (*authorizerv1.FgaGetModelResponse, error) {
	var res authorizerv1.FgaGetModelResponse
	err := c.execute(adminMethodSpec{
		name: "FgaGetModel",
		graphql: &GraphQLRequest{
			Query: `query adminFgaGetModel { _fga_get_model { id dsl } }`,
		},
		graphqlField: "_fga_get_model",
		graphqlWrap:  "model",
		restMethod:   http.MethodGet,
		restPath:     "/v1/admin/fga/model",
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.FgaGetModel(ctx, &authorizerv1.FgaGetModelRequest{})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 26. FgaWriteModel — grpc, rest, gql.
// ---------------------------------------------------------------------------

// FgaWriteModel replaces the fine-grained authorization model. DESTRUCTIVE:
// overwrites the existing authorization model.
func (c *AuthorizerAdminClient) FgaWriteModel(req *authorizerv1.FgaWriteModelRequest) (*authorizerv1.FgaWriteModelResponse, error) {
	var res authorizerv1.FgaWriteModelResponse
	err := c.execute(adminMethodSpec{
		name: "FgaWriteModel",
		graphql: &GraphQLRequest{
			Query:     "mutation fgaWriteModel($data: FgaWriteModelInput!) { _fga_write_model(params: $data) { " + adminFgaModelFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_fga_write_model",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/fga/model",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.FgaWriteModel(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 27. FgaWriteTuples — grpc, rest, gql.
// ---------------------------------------------------------------------------

// FgaWriteTuples writes relationship tuples to the authorization store.
func (c *AuthorizerAdminClient) FgaWriteTuples(req *authorizerv1.FgaWriteTuplesRequest) (*authorizerv1.FgaWriteTuplesResponse, error) {
	var res authorizerv1.FgaWriteTuplesResponse
	err := c.execute(adminMethodSpec{
		name: "FgaWriteTuples",
		graphql: &GraphQLRequest{
			Query:     `mutation fgaWriteTuples($data: FgaWriteTuplesInput!) { _fga_write_tuples(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_fga_write_tuples",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/fga/tuples",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.FgaWriteTuples(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 28. FgaDeleteTuples — grpc, rest, gql.
// ---------------------------------------------------------------------------

// FgaDeleteTuples deletes relationship tuples from the authorization store.
// DESTRUCTIVE: permanently removes the specified tuples.
func (c *AuthorizerAdminClient) FgaDeleteTuples(req *authorizerv1.FgaDeleteTuplesRequest) (*authorizerv1.FgaDeleteTuplesResponse, error) {
	var res authorizerv1.FgaDeleteTuplesResponse
	err := c.execute(adminMethodSpec{
		name: "FgaDeleteTuples",
		graphql: &GraphQLRequest{
			Query:     `mutation fgaDeleteTuples($data: FgaWriteTuplesInput!) { _fga_delete_tuples(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_fga_delete_tuples",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/fga/tuples/delete",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.FgaDeleteTuples(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 29. FgaReadTuples — grpc, rest, gql.
// ---------------------------------------------------------------------------

// FgaReadTuples reads relationship tuples from the authorization store.
func (c *AuthorizerAdminClient) FgaReadTuples(req *authorizerv1.FgaReadTuplesRequest) (*authorizerv1.FgaReadTuplesResponse, error) {
	var res authorizerv1.FgaReadTuplesResponse
	err := c.execute(adminMethodSpec{
		name: "FgaReadTuples",
		graphql: &GraphQLRequest{
			Query:     "query fgaReadTuples($data: FgaReadTuplesInput!) { _fga_read_tuples(params: $data) { tuples { " + adminFgaTupleFields + " } continuation_token } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_fga_read_tuples",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/fga/tuples/read",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.FgaReadTuples(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 30. FgaListUsers — grpc, rest, gql.
// ---------------------------------------------------------------------------

// FgaListUsers lists users that have a given relation to an object.
func (c *AuthorizerAdminClient) FgaListUsers(req *authorizerv1.FgaListUsersRequest) (*authorizerv1.FgaListUsersResponse, error) {
	var res authorizerv1.FgaListUsersResponse
	err := c.execute(adminMethodSpec{
		name: "FgaListUsers",
		graphql: &GraphQLRequest{
			Query:     `query fgaListUsers($data: FgaListUsersInput!) { _fga_list_users(params: $data) { users } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_fga_list_users",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/fga/list_users",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.FgaListUsers(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 31. FgaExpand — grpc, rest, gql.
// ---------------------------------------------------------------------------

// FgaExpand expands the relationship tree for a given object and relation.
func (c *AuthorizerAdminClient) FgaExpand(req *authorizerv1.FgaExpandRequest) (*authorizerv1.FgaExpandResponse, error) {
	var res authorizerv1.FgaExpandResponse
	err := c.execute(adminMethodSpec{
		name: "FgaExpand",
		graphql: &GraphQLRequest{
			Query:     `query fgaExpand($data: FgaExpandInput!) { _fga_expand(params: $data) { tree } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_fga_expand",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/fga/expand",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.FgaExpand(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 32. FgaReset — grpc, rest (no gql).
// ---------------------------------------------------------------------------

// FgaReset deletes the entire fine-grained authorization store. DESTRUCTIVE:
// permanently removes the model and all relationship tuples.
func (c *AuthorizerAdminClient) FgaReset() (*authorizerv1.FgaResetResponse, error) {
	var res authorizerv1.FgaResetResponse
	err := c.execute(adminMethodSpec{
		name: "FgaReset",
		graphql: &GraphQLRequest{
			Query: `mutation adminFgaReset { _fga_reset { message } }`,
		},
		graphqlField: "_fga_reset",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/fga/reset",
		restBody:     &authorizerv1.FgaResetRequest{},
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.FgaReset(ctx, &authorizerv1.FgaResetRequest{})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 33. CreateClient — grpc, rest, gql.
// ---------------------------------------------------------------------------

// CreateClient provisions a new machine/workload identity (service account).
// The client_secret in the response is returned ONCE and can never be
// retrieved again.
func (c *AuthorizerAdminClient) CreateClient(req *authorizerv1.CreateClientRequest) (*authorizerv1.CreateClientResponse, error) {
	var res authorizerv1.CreateClientResponse
	err := c.execute(adminMethodSpec{
		name: "CreateClient",
		graphql: &GraphQLRequest{
			Query:     "mutation createClient($data: CreateClientRequest!) { _create_client(params: $data) { client { " + adminClientFields + " } client_secret } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_create_client",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/create_client",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.CreateClient(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 34. UpdateClient — grpc, rest, gql.
// ---------------------------------------------------------------------------

// UpdateClient updates an existing client's metadata, scopes or active flag.
func (c *AuthorizerAdminClient) UpdateClient(req *authorizerv1.UpdateClientRequest) (*authorizerv1.UpdateClientResponse, error) {
	var res authorizerv1.UpdateClientResponse
	err := c.execute(adminMethodSpec{
		name: "UpdateClient",
		graphql: &GraphQLRequest{
			Query:     "mutation updateClient($data: UpdateClientRequest!) { _update_client(params: $data) { " + adminClientFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_update_client",
		graphqlWrap:  "client",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/update_client",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.UpdateClient(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 35. DeleteClient — grpc, rest, gql.
// ---------------------------------------------------------------------------

// DeleteClient deletes a client. DESTRUCTIVE: permanently removes the client;
// tokens already issued to it stop being honoured.
func (c *AuthorizerAdminClient) DeleteClient(req *authorizerv1.DeleteClientRequest) (*authorizerv1.DeleteClientResponse, error) {
	var res authorizerv1.DeleteClientResponse
	err := c.execute(adminMethodSpec{
		name: "DeleteClient",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteClient($data: ClientRequest!) { _delete_client(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_client",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_client",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteClient(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 36. RotateClientSecret — grpc, rest, gql.
// ---------------------------------------------------------------------------

// RotateClientSecret issues a new secret for a client. The new secret is
// returned ONCE; the old secret keeps validating during the server's grace
// window.
func (c *AuthorizerAdminClient) RotateClientSecret(req *authorizerv1.RotateClientSecretRequest) (*authorizerv1.CreateClientResponse, error) {
	var res authorizerv1.CreateClientResponse
	err := c.execute(adminMethodSpec{
		name: "RotateClientSecret",
		graphql: &GraphQLRequest{
			Query:     "mutation rotateClientSecret($data: ClientRequest!) { _rotate_client_secret(params: $data) { client { " + adminClientFields + " } client_secret } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_rotate_client_secret",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/rotate_client_secret",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.RotateClientSecret(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 37. GetClient — grpc, rest, gql.
// ---------------------------------------------------------------------------

// GetClient returns a single client by id. The client_secret is never returned.
func (c *AuthorizerAdminClient) GetClient(req *authorizerv1.GetClientRequest) (*authorizerv1.GetClientResponse, error) {
	var res authorizerv1.GetClientResponse
	err := c.execute(adminMethodSpec{
		name: "GetClient",
		graphql: &GraphQLRequest{
			Query:     "query client($data: ClientRequest!) { _client(params: $data) { " + adminClientFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_client",
		graphqlWrap:  "client",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/client",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.GetClient(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 38. Clients — grpc, rest, gql.
// ---------------------------------------------------------------------------

// Clients returns a paginated list of clients.
func (c *AuthorizerAdminClient) Clients(req *authorizerv1.ClientsRequest) (*authorizerv1.ClientsResponse, error) {
	gqlData := map[string]interface{}{}
	if req.GetPagination() != nil {
		gqlData["pagination"] = req.GetPagination()
	}

	var res authorizerv1.ClientsResponse
	err := c.execute(adminMethodSpec{
		name: "Clients",
		graphql: &GraphQLRequest{
			Query:     "query clients($data: ListClientsRequest) { _clients(params: $data) { " + adminPaginationFields + " clients { " + adminClientFields + " } } }",
			Variables: map[string]interface{}{"data": gqlData},
		},
		graphqlField: "_clients",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/clients",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.Clients(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 39. AddTrustedIssuer — grpc, rest, gql.
// ---------------------------------------------------------------------------

// AddTrustedIssuer registers an external token issuer (K8s SA, SPIFFE, OIDC)
// that may authenticate as the given service account via JWT-bearer assertions.
func (c *AuthorizerAdminClient) AddTrustedIssuer(req *authorizerv1.AddTrustedIssuerRequest) (*authorizerv1.AddTrustedIssuerResponse, error) {
	var res authorizerv1.AddTrustedIssuerResponse
	err := c.execute(adminMethodSpec{
		name: "AddTrustedIssuer",
		graphql: &GraphQLRequest{
			Query:     "mutation addTrustedIssuer($data: AddTrustedIssuerRequest!) { _add_trusted_issuer(params: $data) { " + adminTrustedIssuerFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_add_trusted_issuer",
		graphqlWrap:  "trusted_issuer",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/add_trusted_issuer",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.AddTrustedIssuer(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 40. UpdateTrustedIssuer — grpc, rest, gql.
// ---------------------------------------------------------------------------

// UpdateTrustedIssuer updates an existing trusted issuer.
func (c *AuthorizerAdminClient) UpdateTrustedIssuer(req *authorizerv1.UpdateTrustedIssuerRequest) (*authorizerv1.UpdateTrustedIssuerResponse, error) {
	var res authorizerv1.UpdateTrustedIssuerResponse
	err := c.execute(adminMethodSpec{
		name: "UpdateTrustedIssuer",
		graphql: &GraphQLRequest{
			Query:     "mutation updateTrustedIssuer($data: UpdateTrustedIssuerRequest!) { _update_trusted_issuer(params: $data) { " + adminTrustedIssuerFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_update_trusted_issuer",
		graphqlWrap:  "trusted_issuer",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/update_trusted_issuer",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.UpdateTrustedIssuer(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 41. DeleteTrustedIssuer — grpc, rest, gql.
// ---------------------------------------------------------------------------

// DeleteTrustedIssuer removes a trusted issuer. DESTRUCTIVE: assertions from
// this issuer stop authenticating immediately.
func (c *AuthorizerAdminClient) DeleteTrustedIssuer(req *authorizerv1.DeleteTrustedIssuerRequest) (*authorizerv1.DeleteTrustedIssuerResponse, error) {
	var res authorizerv1.DeleteTrustedIssuerResponse
	err := c.execute(adminMethodSpec{
		name: "DeleteTrustedIssuer",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteTrustedIssuer($data: TrustedIssuerRequest!) { _delete_trusted_issuer(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_trusted_issuer",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_trusted_issuer",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteTrustedIssuer(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 42. GetTrustedIssuer — grpc, rest, gql.
// ---------------------------------------------------------------------------

// GetTrustedIssuer returns a single trusted issuer by id.
func (c *AuthorizerAdminClient) GetTrustedIssuer(req *authorizerv1.GetTrustedIssuerRequest) (*authorizerv1.GetTrustedIssuerResponse, error) {
	var res authorizerv1.GetTrustedIssuerResponse
	err := c.execute(adminMethodSpec{
		name: "GetTrustedIssuer",
		graphql: &GraphQLRequest{
			Query:     "query trustedIssuer($data: TrustedIssuerRequest!) { _trusted_issuer(params: $data) { " + adminTrustedIssuerFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_trusted_issuer",
		graphqlWrap:  "trusted_issuer",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/trusted_issuer",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.GetTrustedIssuer(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 43. TrustedIssuers — grpc, rest, gql.
// ---------------------------------------------------------------------------

// TrustedIssuers returns a paginated list of trusted issuers, optionally
// filtered by service account.
func (c *AuthorizerAdminClient) TrustedIssuers(req *authorizerv1.TrustedIssuersRequest) (*authorizerv1.TrustedIssuersResponse, error) {
	gqlData := map[string]interface{}{}
	if req.GetPagination() != nil {
		gqlData["pagination"] = req.GetPagination()
	}
	if req.ServiceAccountId != nil {
		gqlData["service_account_id"] = req.GetServiceAccountId()
	}

	var res authorizerv1.TrustedIssuersResponse
	err := c.execute(adminMethodSpec{
		name: "TrustedIssuers",
		graphql: &GraphQLRequest{
			Query:     "query trustedIssuers($data: ListTrustedIssuersRequest) { _trusted_issuers(params: $data) { " + adminPaginationFields + " trusted_issuers { " + adminTrustedIssuerFields + " } } }",
			Variables: map[string]interface{}{"data": gqlData},
		},
		graphqlField: "_trusted_issuers",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/trusted_issuers",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.TrustedIssuers(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 44. CreateSamlServiceProvider — grpc, rest, gql.
// ---------------------------------------------------------------------------

// CreateSamlServiceProvider registers a downstream SAML 2.0 SP that Authorizer
// (acting as the IdP) issues signed assertions to.
func (c *AuthorizerAdminClient) CreateSamlServiceProvider(req *authorizerv1.CreateSamlServiceProviderRequest) (*authorizerv1.CreateSamlServiceProviderResponse, error) {
	var res authorizerv1.CreateSamlServiceProviderResponse
	err := c.execute(adminMethodSpec{
		name: "CreateSamlServiceProvider",
		graphql: &GraphQLRequest{
			Query:     "mutation createSamlServiceProvider($data: CreateSAMLServiceProviderRequest!) { _create_saml_service_provider(params: $data) { " + adminSamlSPFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_create_saml_service_provider",
		graphqlWrap:  "saml_service_provider",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/create_saml_service_provider",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.CreateSamlServiceProvider(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 45. UpdateSamlServiceProvider — grpc, rest, gql.
// ---------------------------------------------------------------------------

// UpdateSamlServiceProvider updates a downstream SP's name, endpoints,
// certificate, attribute mapping, or active state.
func (c *AuthorizerAdminClient) UpdateSamlServiceProvider(req *authorizerv1.UpdateSamlServiceProviderRequest) (*authorizerv1.UpdateSamlServiceProviderResponse, error) {
	var res authorizerv1.UpdateSamlServiceProviderResponse
	err := c.execute(adminMethodSpec{
		name: "UpdateSamlServiceProvider",
		graphql: &GraphQLRequest{
			Query:     "mutation updateSamlServiceProvider($data: UpdateSAMLServiceProviderRequest!) { _update_saml_service_provider(params: $data) { " + adminSamlSPFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_update_saml_service_provider",
		graphqlWrap:  "saml_service_provider",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/update_saml_service_provider",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.UpdateSamlServiceProvider(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 46. DeleteSamlServiceProvider — grpc, rest, gql.
// ---------------------------------------------------------------------------

// DeleteSamlServiceProvider deletes a downstream SP by id. DESTRUCTIVE: SSO
// assertions to this SP stop being issued immediately.
func (c *AuthorizerAdminClient) DeleteSamlServiceProvider(req *authorizerv1.DeleteSamlServiceProviderRequest) (*authorizerv1.DeleteSamlServiceProviderResponse, error) {
	var res authorizerv1.DeleteSamlServiceProviderResponse
	err := c.execute(adminMethodSpec{
		name: "DeleteSamlServiceProvider",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteSamlServiceProvider($data: SAMLServiceProviderRequest!) { _delete_saml_service_provider(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_saml_service_provider",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_saml_service_provider",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteSamlServiceProvider(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 47. GetSamlServiceProvider — grpc, rest, gql.
// ---------------------------------------------------------------------------

// GetSamlServiceProvider returns a single downstream SP by id.
func (c *AuthorizerAdminClient) GetSamlServiceProvider(req *authorizerv1.GetSamlServiceProviderRequest) (*authorizerv1.GetSamlServiceProviderResponse, error) {
	var res authorizerv1.GetSamlServiceProviderResponse
	err := c.execute(adminMethodSpec{
		name: "GetSamlServiceProvider",
		graphql: &GraphQLRequest{
			Query:     "query samlServiceProvider($data: SAMLServiceProviderRequest!) { _saml_service_provider(params: $data) { " + adminSamlSPFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_saml_service_provider",
		graphqlWrap:  "saml_service_provider",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/saml_service_provider",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.GetSamlServiceProvider(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 48. ListSamlServiceProviders — grpc, rest, gql.
// ---------------------------------------------------------------------------

// ListSamlServiceProviders returns a paginated list of downstream SPs for an org.
func (c *AuthorizerAdminClient) ListSamlServiceProviders(req *authorizerv1.ListSamlServiceProvidersRequest) (*authorizerv1.ListSamlServiceProvidersResponse, error) {
	gqlData := map[string]interface{}{"org_id": req.GetOrgId()}
	if req.GetPagination() != nil {
		gqlData["pagination"] = req.GetPagination()
	}

	var res authorizerv1.ListSamlServiceProvidersResponse
	err := c.execute(adminMethodSpec{
		name: "ListSamlServiceProviders",
		graphql: &GraphQLRequest{
			Query:     "query listSamlServiceProviders($data: ListSAMLServiceProvidersRequest!) { _list_saml_service_providers(params: $data) { " + adminPaginationFields + " saml_service_providers { " + adminSamlSPFields + " } } }",
			Variables: map[string]interface{}{"data": gqlData},
		},
		graphqlField: "_list_saml_service_providers",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/saml_service_providers",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.ListSamlServiceProviders(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 49. RotateSamlIdpCert — grpc, rest, gql.
// ---------------------------------------------------------------------------

// RotateSamlIdpCert generates a new current signing keypair for an org's SAML
// IdP, demoting the previous current key.
func (c *AuthorizerAdminClient) RotateSamlIdpCert(req *authorizerv1.RotateSamlIdpCertRequest) (*authorizerv1.RotateSamlIdpCertResponse, error) {
	var res authorizerv1.RotateSamlIdpCertResponse
	err := c.execute(adminMethodSpec{
		name: "RotateSamlIdpCert",
		graphql: &GraphQLRequest{
			Query:     "mutation rotateSamlIdpCert($data: RotateSAMLIDPCertRequest!) { _rotate_saml_idp_cert(params: $data) { " + adminSamlIdpKeyFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_rotate_saml_idp_cert",
		graphqlWrap:  "saml_idp_key",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/rotate_saml_idp_cert",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.RotateSamlIdpCert(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 50. RetireSamlIdpKey — grpc, rest, gql.
// ---------------------------------------------------------------------------

// RetireSamlIdpKey retires a published-but-not-signing SAML IdP key by id.
// DESTRUCTIVE: the key stops being published in IdP metadata.
func (c *AuthorizerAdminClient) RetireSamlIdpKey(req *authorizerv1.RetireSamlIdpKeyRequest) (*authorizerv1.RetireSamlIdpKeyResponse, error) {
	var res authorizerv1.RetireSamlIdpKeyResponse
	err := c.execute(adminMethodSpec{
		name: "RetireSamlIdpKey",
		graphql: &GraphQLRequest{
			Query:     `mutation retireSamlIdpKey($data: RetireSAMLIDPKeyRequest!) { _retire_saml_idp_key(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_retire_saml_idp_key",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/retire_saml_idp_key",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.RetireSamlIdpKey(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 51. ListSamlIdpKeys — grpc, rest, gql.
// ---------------------------------------------------------------------------

// ListSamlIdpKeys returns all SAML IdP signing keys for an org.
func (c *AuthorizerAdminClient) ListSamlIdpKeys(req *authorizerv1.ListSamlIdpKeysRequest) (*authorizerv1.ListSamlIdpKeysResponse, error) {
	var res authorizerv1.ListSamlIdpKeysResponse
	err := c.execute(adminMethodSpec{
		name: "ListSamlIdpKeys",
		graphql: &GraphQLRequest{
			Query:     "query listSamlIdpKeys($data: ListSAMLIDPKeysRequest!) { _list_saml_idp_keys(params: $data) { " + adminSamlIdpKeyFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_list_saml_idp_keys",
		// The GraphQL query returns a bare array; wrap it to match the proto's
		// {saml_idp_keys: [...]} response shape.
		graphqlWrap: "saml_idp_keys",
		restMethod:  http.MethodPost,
		restPath:    "/v1/admin/saml_idp_keys",
		restBody:    req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.ListSamlIdpKeys(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// 52. ImportSamlSpMetadata — grpc, rest, gql.
// ---------------------------------------------------------------------------

// ImportSamlSpMetadata parses pasted SP metadata XML and returns the fields to
// prefill a create call. It does NOT create a record and performs no remote
// fetch.
func (c *AuthorizerAdminClient) ImportSamlSpMetadata(req *authorizerv1.ImportSamlSpMetadataRequest) (*authorizerv1.ImportSamlSpMetadataResponse, error) {
	var res authorizerv1.ImportSamlSpMetadataResponse
	err := c.execute(adminMethodSpec{
		name: "ImportSamlSpMetadata",
		graphql: &GraphQLRequest{
			Query:     `mutation importSamlSpMetadata($data: ImportSAMLSPMetadataRequest!) { _import_saml_sp_metadata(params: $data) { entity_id acs_url certificate } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_import_saml_sp_metadata",
		graphqlWrap:  "result",
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/import_saml_sp_metadata",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.ImportSamlSpMetadata(ctx, req)
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// GraphQL-only admin operations. These have no gRPC stub or REST endpoint, so
// only the graphql protocol is supported. The proto definitions do not include
// these operations, so local request/response types are declared here.
// ---------------------------------------------------------------------------

// AdminSignupRequest is the request for AdminSignup.
type AdminSignupRequest struct {
	AdminSecret string `json:"admin_secret"`
}

// AdminSignup sets the admin secret for a fresh instance (gql only).
func (c *AuthorizerAdminClient) AdminSignup(req *AdminSignupRequest) (*Response, error) {
	var res Response
	err := c.execute(adminMethodSpec{
		name: "AdminSignup",
		graphql: &GraphQLRequest{
			Query:     `mutation adminSignup($data: AdminSignupRequest!) { _admin_signup(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_admin_signup",
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateEnvRequest is a flexible map of environment keys to update.
type UpdateEnvRequest map[string]interface{}

// UpdateEnv updates the instance environment configuration (gql only).
func (c *AuthorizerAdminClient) UpdateEnv(req *UpdateEnvRequest) (*Response, error) {
	var res Response
	err := c.execute(adminMethodSpec{
		name: "UpdateEnv",
		graphql: &GraphQLRequest{
			Query:     `mutation updateEnv($data: UpdateEnvInput!) { _update_env(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_update_env",
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GenerateJWTKeysRequest is the request for GenerateJWTKeys.
type GenerateJWTKeysRequest struct {
	Type string `json:"type"`
}

// GenerateJWTKeysResponse is the response for GenerateJWTKeys.
type GenerateJWTKeysResponse struct {
	Secret     *string `json:"secret"`
	PublicKey  *string `json:"public_key"`
	PrivateKey *string `json:"private_key"`
}

// GenerateJWTKeys generates a new set of JWT signing keys (gql only).
func (c *AuthorizerAdminClient) GenerateJWTKeys(req *GenerateJWTKeysRequest) (*GenerateJWTKeysResponse, error) {
	var res GenerateJWTKeysResponse
	err := c.execute(adminMethodSpec{
		name: "GenerateJWTKeys",
		graphql: &GraphQLRequest{
			Query:     `query generateJwtKeys($data: GenerateJWTKeysRequest!) { _generate_jwt_keys(params: $data) { secret public_key private_key } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_generate_jwt_keys",
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// Organizations, org members, org OIDC/SAML connections and SCIM endpoints are
// GraphQL-only: the admin proto has no RPCs for them, so rest/grpc return a
// clear unsupported-protocol error. Types mirror the GraphQL schema.
// ---------------------------------------------------------------------------

// PaginationRequest mirrors the GraphQL PaginationRequest input.
type PaginationRequest struct {
	Limit int64 `json:"limit,omitempty"`
	Page  int64 `json:"page,omitempty"`
}

// Pagination mirrors the GraphQL Pagination response type.
type Pagination struct {
	Limit  int64 `json:"limit"`
	Page   int64 `json:"page"`
	Offset int64 `json:"offset"`
	Total  int64 `json:"total"`
}

// Organization defines attributes of an organization.
type Organization struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	DisplayName *string `json:"display_name"`
	Enabled     bool    `json:"enabled"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

// Organizations is a paginated list of organizations.
type Organizations struct {
	Pagination    *Pagination     `json:"pagination"`
	Organizations []*Organization `json:"organizations"`
}

// OrgMember defines a user's membership in an organization.
type OrgMember struct {
	ID         string   `json:"id"`
	OrgID      string   `json:"org_id"`
	UserID     string   `json:"user_id"`
	Email      *string  `json:"email"`
	GivenName  *string  `json:"given_name"`
	FamilyName *string  `json:"family_name"`
	Roles      []string `json:"roles"`
	CreatedAt  int64    `json:"created_at"`
	UpdatedAt  int64    `json:"updated_at"`
}

// OrgMembers is a paginated list of organization members.
type OrgMembers struct {
	Pagination *Pagination  `json:"pagination"`
	OrgMembers []*OrgMember `json:"org_members"`
}

// CreateOrganizationRequest is the request for CreateOrganization. Name must
// be a unique, URL-safe slug.
type CreateOrganizationRequest struct {
	Name        string  `json:"name"`
	DisplayName *string `json:"display_name,omitempty"`
}

// UpdateOrganizationRequest is the request for UpdateOrganization.
type UpdateOrganizationRequest struct {
	ID          string  `json:"id"`
	Name        *string `json:"name,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
}

// OrganizationRequest identifies an organization by id.
type OrganizationRequest struct {
	ID string `json:"id"`
}

// ListOrganizationsRequest is the request for Organizations.
type ListOrganizationsRequest struct {
	Pagination *PaginationRequest `json:"pagination,omitempty"`
}

// AddOrgMemberRequest is the request for AddOrgMember. Roles defaults to an
// empty set when omitted.
type AddOrgMemberRequest struct {
	OrgID  string   `json:"org_id"`
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles,omitempty"`
}

// RemoveOrgMemberRequest is the request for RemoveOrgMember.
type RemoveOrgMemberRequest struct {
	OrgID  string `json:"org_id"`
	UserID string `json:"user_id"`
}

// ListOrgMembersRequest is the request for OrgMembers.
type ListOrgMembersRequest struct {
	OrgID      string             `json:"org_id"`
	Pagination *PaginationRequest `json:"pagination,omitempty"`
}

// CreateOrganization creates an organization (gql only).
func (c *AuthorizerAdminClient) CreateOrganization(req *CreateOrganizationRequest) (*Organization, error) {
	var res Organization
	err := c.execute(adminMethodSpec{
		name: "CreateOrganization",
		graphql: &GraphQLRequest{
			Query:     "mutation createOrganization($data: CreateOrganizationRequest!) { _create_organization(params: $data) { " + adminOrgFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_create_organization",
		responseUnwrap: "organization",
		restResponse:   func() proto.Message { return &authorizerv1.CreateOrganizationResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/create_organization",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.CreateOrganization(ctx, &authorizerv1.CreateOrganizationRequest{
				Name: req.GetName(), DisplayName: req.GetDisplayName(),
			})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateOrganization updates an organization (gql only).
func (c *AuthorizerAdminClient) UpdateOrganization(req *UpdateOrganizationRequest) (*Organization, error) {
	var res Organization
	err := c.execute(adminMethodSpec{
		name: "UpdateOrganization",
		graphql: &GraphQLRequest{
			Query:     "mutation updateOrganization($data: UpdateOrganizationRequest!) { _update_organization(params: $data) { " + adminOrgFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_update_organization",
		responseUnwrap: "organization",
		restResponse:   func() proto.Message { return &authorizerv1.UpdateOrganizationResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/update_organization",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.UpdateOrganization(ctx, &authorizerv1.UpdateOrganizationRequest{
				Id: req.GetID(), Name: req.GetName(), DisplayName: req.GetDisplayName(), Enabled: req.GetEnabled(),
			})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteOrganization deletes an organization (gql only). DESTRUCTIVE:
// permanently removes the organization and its memberships/connections.
func (c *AuthorizerAdminClient) DeleteOrganization(req *OrganizationRequest) (*Response, error) {
	var res Response
	err := c.execute(adminMethodSpec{
		name: "DeleteOrganization",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteOrganization($data: OrganizationRequest!) { _delete_organization(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_organization",
		restResponse: func() proto.Message { return &authorizerv1.DeleteOrganizationResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_organization",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteOrganization(ctx, &authorizerv1.DeleteOrganizationRequest{Id: req.GetID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// AddOrgMember adds a user to an organization (gql only).
func (c *AuthorizerAdminClient) AddOrgMember(req *AddOrgMemberRequest) (*OrgMember, error) {
	var res OrgMember
	err := c.execute(adminMethodSpec{
		name: "AddOrgMember",
		graphql: &GraphQLRequest{
			Query:     "mutation addOrgMember($data: AddOrgMemberRequest!) { _add_org_member(params: $data) { " + adminOrgMemberFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_add_org_member",
		responseUnwrap: "org_member",
		restResponse:   func() proto.Message { return &authorizerv1.AddOrgMemberResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/add_org_member",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.AddOrgMember(ctx, &authorizerv1.AddOrgMemberRequest{
				OrgId: req.GetOrgID(), UserId: req.GetUserID(), Roles: req.GetRoles(),
			})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// RemoveOrgMember removes a user from an organization (gql only).
func (c *AuthorizerAdminClient) RemoveOrgMember(req *RemoveOrgMemberRequest) (*Response, error) {
	var res Response
	err := c.execute(adminMethodSpec{
		name: "RemoveOrgMember",
		graphql: &GraphQLRequest{
			Query:     `mutation removeOrgMember($data: RemoveOrgMemberRequest!) { _remove_org_member(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_remove_org_member",
		restResponse: func() proto.Message { return &authorizerv1.RemoveOrgMemberResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/remove_org_member",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.RemoveOrgMember(ctx, &authorizerv1.RemoveOrgMemberRequest{OrgId: req.GetOrgID(), UserId: req.GetUserID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetOrganization returns a single organization by id (gql only).
func (c *AuthorizerAdminClient) GetOrganization(req *OrganizationRequest) (*Organization, error) {
	var res Organization
	err := c.execute(adminMethodSpec{
		name: "GetOrganization",
		graphql: &GraphQLRequest{
			Query:     "query organization($data: OrganizationRequest!) { _organization(params: $data) { " + adminOrgFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_organization",
		responseUnwrap: "organization",
		restResponse:   func() proto.Message { return &authorizerv1.GetOrganizationResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/organization",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.GetOrganization(ctx, &authorizerv1.GetOrganizationRequest{Id: req.GetID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Organizations returns a paginated list of organizations (gql only).
func (c *AuthorizerAdminClient) Organizations(req *ListOrganizationsRequest) (*Organizations, error) {
	var res Organizations
	err := c.execute(adminMethodSpec{
		name: "Organizations",
		graphql: &GraphQLRequest{
			Query:     "query organizations($data: ListOrganizationsRequest) { _organizations(params: $data) { " + adminPaginationFields + " organizations { " + adminOrgFields + " } } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_organizations",
		restResponse: func() proto.Message { return &authorizerv1.OrganizationsResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/organizations",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.Organizations(ctx, &authorizerv1.OrganizationsRequest{Pagination: req.protoPagination()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// OrgMembers returns a paginated list of an organization's members (gql only).
func (c *AuthorizerAdminClient) OrgMembers(req *ListOrgMembersRequest) (*OrgMembers, error) {
	var res OrgMembers
	err := c.execute(adminMethodSpec{
		name: "OrgMembers",
		graphql: &GraphQLRequest{
			Query:     "query orgMembers($data: ListOrgMembersRequest!) { _org_members(params: $data) { " + adminPaginationFields + " org_members { " + adminOrgMemberFields + " } } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_org_members",
		restResponse: func() proto.Message { return &authorizerv1.OrgMembersResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/org_members",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.OrgMembers(ctx, &authorizerv1.OrgMembersRequest{OrgId: req.GetOrgID(), Pagination: req.protoPagination()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// OrgOIDCConnection defines an organization's upstream OIDC SSO connection.
type OrgOIDCConnection struct {
	ID          string  `json:"id"`
	OrgID       string  `json:"org_id"`
	Name        string  `json:"name"`
	IssuerURL   string  `json:"issuer_url"`
	SSOClientID string  `json:"sso_client_id"`
	Scopes      *string `json:"scopes"`
	RedirectURI *string `json:"redirect_uri"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

// CreateOrgOIDCConnectionRequest is the request for CreateOrgOIDCConnection.
// ClientID/ClientSecret are the credentials Authorizer holds at the upstream
// IdP; the secret is stored encrypted and never returned.
type CreateOrgOIDCConnectionRequest struct {
	OrgID        string  `json:"org_id"`
	Name         string  `json:"name"`
	IssuerURL    string  `json:"issuer_url"`
	ClientID     string  `json:"client_id"`
	ClientSecret string  `json:"client_secret"`
	Scopes       *string `json:"scopes,omitempty"`
	RedirectURI  *string `json:"redirect_uri,omitempty"`
}

// UpdateOrgOIDCConnectionRequest is the request for UpdateOrgOIDCConnection.
// Supplying ClientSecret rotates it; omitting leaves the stored secret intact.
type UpdateOrgOIDCConnectionRequest struct {
	ID           string  `json:"id"`
	Name         *string `json:"name,omitempty"`
	IssuerURL    *string `json:"issuer_url,omitempty"`
	ClientID     *string `json:"client_id,omitempty"`
	ClientSecret *string `json:"client_secret,omitempty"`
	Scopes       *string `json:"scopes,omitempty"`
	RedirectURI  *string `json:"redirect_uri,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// OrgOIDCConnectionRequest looks a connection up by id OR by org id (supply
// exactly one).
type OrgOIDCConnectionRequest struct {
	ID    *string `json:"id,omitempty"`
	OrgID *string `json:"org_id,omitempty"`
}

// CreateOrgOIDCConnection creates an org OIDC SSO connection (gql only).
func (c *AuthorizerAdminClient) CreateOrgOIDCConnection(req *CreateOrgOIDCConnectionRequest) (*OrgOIDCConnection, error) {
	var res OrgOIDCConnection
	err := c.execute(adminMethodSpec{
		name: "CreateOrgOIDCConnection",
		graphql: &GraphQLRequest{
			Query:     "mutation createOrgOidcConnection($data: CreateOrgOIDCConnectionRequest!) { _create_org_oidc_connection(params: $data) { " + adminOrgOIDCConnFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_create_org_oidc_connection",
		responseUnwrap: "org_oidc_connection",
		restResponse:   func() proto.Message { return &authorizerv1.CreateOrgOidcConnectionResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/create_org_oidc_connection",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.CreateOrgOidcConnection(ctx, &authorizerv1.CreateOrgOidcConnectionRequest{
				OrgId: req.GetOrgID(), Name: req.GetName(), IssuerUrl: req.GetIssuerURL(),
				ClientId: req.GetClientID(), ClientSecret: req.GetClientSecret(),
				Scopes: req.GetScopes(), RedirectUri: req.GetRedirectURI(),
			})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateOrgOIDCConnection updates an org OIDC SSO connection (gql only).
func (c *AuthorizerAdminClient) UpdateOrgOIDCConnection(req *UpdateOrgOIDCConnectionRequest) (*OrgOIDCConnection, error) {
	var res OrgOIDCConnection
	err := c.execute(adminMethodSpec{
		name: "UpdateOrgOIDCConnection",
		graphql: &GraphQLRequest{
			Query:     "mutation updateOrgOidcConnection($data: UpdateOrgOIDCConnectionRequest!) { _update_org_oidc_connection(params: $data) { " + adminOrgOIDCConnFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_update_org_oidc_connection",
		responseUnwrap: "org_oidc_connection",
		restResponse:   func() proto.Message { return &authorizerv1.UpdateOrgOidcConnectionResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/update_org_oidc_connection",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.UpdateOrgOidcConnection(ctx, &authorizerv1.UpdateOrgOidcConnectionRequest{
				Id: req.GetID(), Name: req.GetName(), IssuerUrl: req.GetIssuerURL(),
				ClientId: req.GetClientID(), ClientSecret: req.GetClientSecret(),
				Scopes: req.GetScopes(), RedirectUri: req.GetRedirectURI(), IsActive: req.GetIsActive(),
			})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteOrgOIDCConnection deletes an org OIDC SSO connection (gql only).
// DESTRUCTIVE: SSO logins through this connection stop working immediately.
func (c *AuthorizerAdminClient) DeleteOrgOIDCConnection(req *OrgOIDCConnectionRequest) (*Response, error) {
	var res Response
	err := c.execute(adminMethodSpec{
		name: "DeleteOrgOIDCConnection",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteOrgOidcConnection($data: OrgOIDCConnectionRequest!) { _delete_org_oidc_connection(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_org_oidc_connection",
		restResponse: func() proto.Message { return &authorizerv1.DeleteOrgOidcConnectionResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_org_oidc_connection",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteOrgOidcConnection(ctx, &authorizerv1.DeleteOrgOidcConnectionRequest{Id: req.GetID(), OrgId: req.GetOrgID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetOrgOIDCConnection returns an org OIDC SSO connection by id or org id
// (gql only).
func (c *AuthorizerAdminClient) GetOrgOIDCConnection(req *OrgOIDCConnectionRequest) (*OrgOIDCConnection, error) {
	var res OrgOIDCConnection
	err := c.execute(adminMethodSpec{
		name: "GetOrgOIDCConnection",
		graphql: &GraphQLRequest{
			Query:     "query orgOidcConnection($data: OrgOIDCConnectionRequest!) { _org_oidc_connection(params: $data) { " + adminOrgOIDCConnFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_org_oidc_connection",
		responseUnwrap: "org_oidc_connection",
		restResponse:   func() proto.Message { return &authorizerv1.GetOrgOidcConnectionResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/org_oidc_connection",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.GetOrgOidcConnection(ctx, &authorizerv1.GetOrgOidcConnectionRequest{Id: req.GetID(), OrgId: req.GetOrgID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// OrgSAMLConnection defines an organization's upstream SAML SSO connection.
type OrgSAMLConnection struct {
	ID                string  `json:"id"`
	OrgID             string  `json:"org_id"`
	Name              string  `json:"name"`
	IdpEntityID       string  `json:"idp_entity_id"`
	IdpSSOURL         *string `json:"idp_sso_url"`
	SpEntityID        *string `json:"sp_entity_id"`
	AcsURL            *string `json:"acs_url"`
	AttributeMapping  *string `json:"attribute_mapping"`
	AllowIdpInitiated bool    `json:"allow_idp_initiated"`
	IsActive          bool    `json:"is_active"`
	CreatedAt         int64   `json:"created_at"`
	UpdatedAt         int64   `json:"updated_at"`
}

// CreateOrgSAMLConnectionRequest is the request for CreateOrgSAMLConnection.
// IdpCertificate is the IdP X.509 signing certificate (PEM); assertion
// signatures are validated only against it.
type CreateOrgSAMLConnectionRequest struct {
	OrgID             string  `json:"org_id"`
	Name              string  `json:"name"`
	IdpEntityID       string  `json:"idp_entity_id"`
	IdpSSOURL         string  `json:"idp_sso_url"`
	IdpCertificate    string  `json:"idp_certificate"`
	SpEntityID        *string `json:"sp_entity_id,omitempty"`
	AcsURL            *string `json:"acs_url,omitempty"`
	AttributeMapping  *string `json:"attribute_mapping,omitempty"`
	AllowIdpInitiated *bool   `json:"allow_idp_initiated,omitempty"`
}

// UpdateOrgSAMLConnectionRequest is the request for UpdateOrgSAMLConnection.
// Supplying IdpCertificate replaces it; omitting leaves the stored cert intact.
type UpdateOrgSAMLConnectionRequest struct {
	ID                string  `json:"id"`
	Name              *string `json:"name,omitempty"`
	IdpEntityID       *string `json:"idp_entity_id,omitempty"`
	IdpSSOURL         *string `json:"idp_sso_url,omitempty"`
	IdpCertificate    *string `json:"idp_certificate,omitempty"`
	SpEntityID        *string `json:"sp_entity_id,omitempty"`
	AcsURL            *string `json:"acs_url,omitempty"`
	AttributeMapping  *string `json:"attribute_mapping,omitempty"`
	AllowIdpInitiated *bool   `json:"allow_idp_initiated,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
}

// OrgSAMLConnectionRequest looks a connection up by id OR by org id (supply
// exactly one).
type OrgSAMLConnectionRequest struct {
	ID    *string `json:"id,omitempty"`
	OrgID *string `json:"org_id,omitempty"`
}

// CreateOrgSAMLConnection creates an org SAML SSO connection (gql only).
func (c *AuthorizerAdminClient) CreateOrgSAMLConnection(req *CreateOrgSAMLConnectionRequest) (*OrgSAMLConnection, error) {
	var res OrgSAMLConnection
	err := c.execute(adminMethodSpec{
		name: "CreateOrgSAMLConnection",
		graphql: &GraphQLRequest{
			Query:     "mutation createOrgSamlConnection($data: CreateOrgSAMLConnectionRequest!) { _create_org_saml_connection(params: $data) { " + adminOrgSAMLConnFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_create_org_saml_connection",
		responseUnwrap: "org_saml_connection",
		restResponse:   func() proto.Message { return &authorizerv1.CreateOrgSamlConnectionResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/create_org_saml_connection",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.CreateOrgSamlConnection(ctx, &authorizerv1.CreateOrgSamlConnectionRequest{
				OrgId: req.GetOrgID(), Name: req.GetName(), IdpEntityId: req.GetIdpEntityID(),
				IdpSsoUrl: req.GetIdpSSOURL(), IdpCertificate: req.GetIdpCertificate(),
				SpEntityId: req.GetSpEntityID(), AcsUrl: req.GetAcsURL(),
				AttributeMapping: req.GetAttributeMapping(), AllowIdpInitiated: req.GetAllowIdpInitiated(),
			})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateOrgSAMLConnection updates an org SAML SSO connection (gql only).
func (c *AuthorizerAdminClient) UpdateOrgSAMLConnection(req *UpdateOrgSAMLConnectionRequest) (*OrgSAMLConnection, error) {
	var res OrgSAMLConnection
	err := c.execute(adminMethodSpec{
		name: "UpdateOrgSAMLConnection",
		graphql: &GraphQLRequest{
			Query:     "mutation updateOrgSamlConnection($data: UpdateOrgSAMLConnectionRequest!) { _update_org_saml_connection(params: $data) { " + adminOrgSAMLConnFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_update_org_saml_connection",
		responseUnwrap: "org_saml_connection",
		restResponse:   func() proto.Message { return &authorizerv1.UpdateOrgSamlConnectionResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/update_org_saml_connection",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.UpdateOrgSamlConnection(ctx, &authorizerv1.UpdateOrgSamlConnectionRequest{
				Id: req.GetID(), Name: req.GetName(), IdpEntityId: req.GetIdpEntityID(),
				IdpSsoUrl: req.GetIdpSSOURL(), IdpCertificate: req.GetIdpCertificate(),
				SpEntityId: req.GetSpEntityID(), AcsUrl: req.GetAcsURL(),
				AttributeMapping: req.GetAttributeMapping(), AllowIdpInitiated: req.GetAllowIdpInitiated(),
				IsActive: req.GetIsActive(),
			})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteOrgSAMLConnection deletes an org SAML SSO connection (gql only).
// DESTRUCTIVE: SSO logins through this connection stop working immediately.
func (c *AuthorizerAdminClient) DeleteOrgSAMLConnection(req *OrgSAMLConnectionRequest) (*Response, error) {
	var res Response
	err := c.execute(adminMethodSpec{
		name: "DeleteOrgSAMLConnection",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteOrgSamlConnection($data: OrgSAMLConnectionRequest!) { _delete_org_saml_connection(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_org_saml_connection",
		restResponse: func() proto.Message { return &authorizerv1.DeleteOrgSamlConnectionResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_org_saml_connection",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteOrgSamlConnection(ctx, &authorizerv1.DeleteOrgSamlConnectionRequest{Id: req.GetID(), OrgId: req.GetOrgID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetOrgSAMLConnection returns an org SAML SSO connection by id or org id
// (gql only).
func (c *AuthorizerAdminClient) GetOrgSAMLConnection(req *OrgSAMLConnectionRequest) (*OrgSAMLConnection, error) {
	var res OrgSAMLConnection
	err := c.execute(adminMethodSpec{
		name: "GetOrgSAMLConnection",
		graphql: &GraphQLRequest{
			Query:     "query orgSamlConnection($data: OrgSAMLConnectionRequest!) { _org_saml_connection(params: $data) { " + adminOrgSAMLConnFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_org_saml_connection",
		responseUnwrap: "org_saml_connection",
		restResponse:   func() proto.Message { return &authorizerv1.GetOrgSamlConnectionResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/org_saml_connection",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.GetOrgSamlConnection(ctx, &authorizerv1.GetOrgSamlConnectionRequest{Id: req.GetID(), OrgId: req.GetOrgID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ScimEndpoint defines an organization's SCIM provisioning endpoint.
type ScimEndpoint struct {
	ID        string `json:"id"`
	OrgID     string `json:"org_id"`
	Enabled   bool   `json:"enabled"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// CreateScimEndpointResponse carries the SCIM endpoint plus its bearer token.
// The token is returned ONCE at creation and ONCE at rotation; it can never be
// retrieved again.
type CreateScimEndpointResponse struct {
	ScimEndpoint *ScimEndpoint `json:"scim_endpoint"`
	Token        string        `json:"token"`
}

// CreateScimEndpointRequest is the request for CreateScimEndpoint.
type CreateScimEndpointRequest struct {
	OrgID string `json:"org_id"`
}

// ScimEndpointRequest identifies an organization's SCIM endpoint by org id.
type ScimEndpointRequest struct {
	OrgID string `json:"org_id"`
}

// scimEndpointResponseFragment selects the SCIM endpoint + one-time token.
const scimEndpointResponseFragment = "scim_endpoint { " + adminScimEndpointFields + " } token"

// CreateScimEndpoint provisions a SCIM endpoint for an organization (gql only).
func (c *AuthorizerAdminClient) CreateScimEndpoint(req *CreateScimEndpointRequest) (*CreateScimEndpointResponse, error) {
	var res CreateScimEndpointResponse
	err := c.execute(adminMethodSpec{
		name: "CreateScimEndpoint",
		graphql: &GraphQLRequest{
			Query:     "mutation createScimEndpoint($data: CreateScimEndpointRequest!) { _create_scim_endpoint(params: $data) { " + scimEndpointResponseFragment + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_create_scim_endpoint",
		restResponse: func() proto.Message { return &authorizerv1.CreateScimEndpointResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/create_scim_endpoint",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.CreateScimEndpoint(ctx, &authorizerv1.CreateScimEndpointRequest{OrgId: req.GetOrgID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// RotateScimToken rotates the SCIM endpoint's bearer token (gql only). The new
// token is returned ONCE; the old token stops validating.
func (c *AuthorizerAdminClient) RotateScimToken(req *ScimEndpointRequest) (*CreateScimEndpointResponse, error) {
	var res CreateScimEndpointResponse
	err := c.execute(adminMethodSpec{
		name: "RotateScimToken",
		graphql: &GraphQLRequest{
			Query:     "mutation rotateScimToken($data: ScimEndpointRequest!) { _rotate_scim_token(params: $data) { " + scimEndpointResponseFragment + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_rotate_scim_token",
		restResponse: func() proto.Message { return &authorizerv1.CreateScimEndpointResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/rotate_scim_token",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.RotateScimToken(ctx, &authorizerv1.RotateScimTokenRequest{OrgId: req.GetOrgID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteScimEndpoint deletes an organization's SCIM endpoint (gql only).
// DESTRUCTIVE: the IdP's provisioning token stops working immediately.
func (c *AuthorizerAdminClient) DeleteScimEndpoint(req *ScimEndpointRequest) (*Response, error) {
	var res Response
	err := c.execute(adminMethodSpec{
		name: "DeleteScimEndpoint",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteScimEndpoint($data: ScimEndpointRequest!) { _delete_scim_endpoint(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_scim_endpoint",
		restResponse: func() proto.Message { return &authorizerv1.DeleteScimEndpointResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_scim_endpoint",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteScimEndpoint(ctx, &authorizerv1.DeleteScimEndpointRequest{OrgId: req.GetOrgID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetScimEndpoint returns an organization's SCIM endpoint (gql only). The
// bearer token is never returned.
func (c *AuthorizerAdminClient) GetScimEndpoint(req *ScimEndpointRequest) (*ScimEndpoint, error) {
	var res ScimEndpoint
	err := c.execute(adminMethodSpec{
		name: "GetScimEndpoint",
		graphql: &GraphQLRequest{
			Query:     "query scimEndpoint($data: ScimEndpointRequest!) { _scim_endpoint(params: $data) { " + adminScimEndpointFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_scim_endpoint",
		responseUnwrap: "scim_endpoint",
		restResponse:   func() proto.Message { return &authorizerv1.GetScimEndpointResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/scim_endpoint",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.GetScimEndpoint(ctx, &authorizerv1.GetScimEndpointRequest{OrgId: req.GetOrgID()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// UserOrganizations — a user's organizations plus their per-org roles. The
// admin proto has no RPC for this, so only graphql is supported.
// ---------------------------------------------------------------------------

// UserOrganization pairs an organization with the roles a user holds in it.
type UserOrganization struct {
	Organization *Organization `json:"organization"`
	Roles        []string      `json:"roles"`
}

// UserOrganizations is a paginated list of a user's organizations.
type UserOrganizations struct {
	Pagination        *Pagination         `json:"pagination"`
	UserOrganizations []*UserOrganization `json:"user_organizations"`
}

// UserOrganizationsRequest is the request for UserOrganizations.
type UserOrganizationsRequest struct {
	UserID     string             `json:"user_id"`
	Pagination *PaginationRequest `json:"pagination,omitempty"`
}

const adminUserOrgFields = `organization { ` + adminOrgFields + ` } roles`

// UserOrganizations returns the organizations a user belongs to, with the
// roles held per org (gql only).
func (c *AuthorizerAdminClient) UserOrganizations(req *UserOrganizationsRequest) (*UserOrganizations, error) {
	var res UserOrganizations
	err := c.execute(adminMethodSpec{
		name: "UserOrganizations",
		graphql: &GraphQLRequest{
			Query:     "query userOrganizations($data: UserOrganizationsRequest!) { _user_organizations(params: $data) { " + adminPaginationFields + " user_organizations { " + adminUserOrgFields + " } } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_user_organizations",
		restResponse: func() proto.Message { return &authorizerv1.UserOrganizationsResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/user_organizations",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.UserOrganizations(ctx, &authorizerv1.UserOrganizationsRequest{UserId: req.GetUserID(), Pagination: req.protoPagination()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------------------------------------------------------------------------
// Org domains (home-realm discovery). The admin proto has no RPCs for these,
// so only graphql is supported. request/verify/delete are org-admin gated;
// AddVerifiedOrgDomain is super-admin only (trusted-assert, bypasses the DNS
// TXT challenge).
// ---------------------------------------------------------------------------

// OrgDomain is a verified mapping from a DNS domain to exactly one organization,
// used for home-realm discovery.
type OrgDomain struct {
	Domain     string `json:"domain"`
	OrgID      string `json:"org_id"`
	VerifiedAt *int64 `json:"verified_at"`
	CreatedAt  *int64 `json:"created_at"`
	UpdatedAt  *int64 `json:"updated_at"`
}

// OrgDomains is a paginated list of an organization's verified domains.
type OrgDomains struct {
	Pagination *Pagination  `json:"pagination"`
	OrgDomains []*OrgDomain `json:"org_domains"`
}

// OrgDomainChallenge is the DNS TXT record a tenant must publish to prove
// control of a domain. Returned by RequestOrgDomain; no durable row exists
// until the domain is verified.
type OrgDomainChallenge struct {
	Domain      string `json:"domain"`
	RecordType  string `json:"record_type"`
	RecordName  string `json:"record_name"`
	RecordValue string `json:"record_value"`
}

// RequestOrgDomainRequest is the request for RequestOrgDomain.
type RequestOrgDomainRequest struct {
	OrgID  string `json:"org_id"`
	Domain string `json:"domain"`
}

// VerifyOrgDomainRequest is the request for VerifyOrgDomain.
type VerifyOrgDomainRequest struct {
	OrgID  string `json:"org_id"`
	Domain string `json:"domain"`
}

// AddVerifiedOrgDomainRequest is the request for AddVerifiedOrgDomain.
type AddVerifiedOrgDomainRequest struct {
	OrgID  string `json:"org_id"`
	Domain string `json:"domain"`
}

// ListOrgDomainsRequest is the request for OrgDomains.
type ListOrgDomainsRequest struct {
	OrgID      string             `json:"org_id"`
	Pagination *PaginationRequest `json:"pagination,omitempty"`
}

// DeleteOrgDomainRequest is the request for DeleteOrgDomain.
type DeleteOrgDomainRequest struct {
	Domain string `json:"domain"`
}

const adminOrgDomainFields = `domain org_id verified_at created_at updated_at`

// RequestOrgDomain starts domain verification for home-realm discovery,
// returning the DNS TXT record the tenant must publish to prove control of
// the domain (gql only).
func (c *AuthorizerAdminClient) RequestOrgDomain(req *RequestOrgDomainRequest) (*OrgDomainChallenge, error) {
	var res OrgDomainChallenge
	err := c.execute(adminMethodSpec{
		name: "RequestOrgDomain",
		graphql: &GraphQLRequest{
			Query:     `mutation requestOrgDomain($data: RequestOrgDomainRequest!) { _request_org_domain(params: $data) { domain record_type record_name record_value } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_request_org_domain",
		responseUnwrap: "challenge",
		restResponse:   func() proto.Message { return &authorizerv1.RequestOrgDomainResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/request_org_domain",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.RequestOrgDomain(ctx, &authorizerv1.RequestOrgDomainRequest{OrgId: req.GetOrgID(), Domain: req.GetDomain()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// VerifyOrgDomain checks the DNS TXT challenge and, if satisfied, verifies the
// domain (gql only).
func (c *AuthorizerAdminClient) VerifyOrgDomain(req *VerifyOrgDomainRequest) (*OrgDomain, error) {
	var res OrgDomain
	err := c.execute(adminMethodSpec{
		name: "VerifyOrgDomain",
		graphql: &GraphQLRequest{
			Query:     "mutation verifyOrgDomain($data: VerifyOrgDomainRequest!) { _verify_org_domain(params: $data) { " + adminOrgDomainFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_verify_org_domain",
		responseUnwrap: "org_domain",
		restResponse:   func() proto.Message { return &authorizerv1.VerifyOrgDomainResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/verify_org_domain",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.VerifyOrgDomain(ctx, &authorizerv1.VerifyOrgDomainRequest{OrgId: req.GetOrgID(), Domain: req.GetDomain()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// AddVerifiedOrgDomain directly registers a verified domain, bypassing the DNS
// TXT challenge. Super-admin only (gql only).
func (c *AuthorizerAdminClient) AddVerifiedOrgDomain(req *AddVerifiedOrgDomainRequest) (*OrgDomain, error) {
	var res OrgDomain
	err := c.execute(adminMethodSpec{
		name: "AddVerifiedOrgDomain",
		graphql: &GraphQLRequest{
			Query:     "mutation addVerifiedOrgDomain($data: AddVerifiedOrgDomainRequest!) { _add_verified_org_domain(params: $data) { " + adminOrgDomainFields + " } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField:   "_add_verified_org_domain",
		responseUnwrap: "org_domain",
		restResponse:   func() proto.Message { return &authorizerv1.AddVerifiedOrgDomainResponse{} },
		restMethod:     http.MethodPost,
		restPath:       "/v1/admin/add_verified_org_domain",
		restBody:       req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.AddVerifiedOrgDomain(ctx, &authorizerv1.AddVerifiedOrgDomainRequest{OrgId: req.GetOrgID(), Domain: req.GetDomain()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteOrgDomain removes a verified domain (gql only). DESTRUCTIVE: logins
// relying on this domain for home-realm discovery stop resolving to the org.
func (c *AuthorizerAdminClient) DeleteOrgDomain(req *DeleteOrgDomainRequest) (*Response, error) {
	var res Response
	err := c.execute(adminMethodSpec{
		name: "DeleteOrgDomain",
		graphql: &GraphQLRequest{
			Query:     `mutation deleteOrgDomain($data: DeleteOrgDomainRequest!) { _delete_org_domain(params: $data) { message } }`,
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_delete_org_domain",
		restResponse: func() proto.Message { return &authorizerv1.DeleteOrgDomainResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/delete_org_domain",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.DeleteOrgDomain(ctx, &authorizerv1.DeleteOrgDomainRequest{Domain: req.GetDomain()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// OrgDomains returns an organization's verified domains (gql only).
func (c *AuthorizerAdminClient) OrgDomains(req *ListOrgDomainsRequest) (*OrgDomains, error) {
	var res OrgDomains
	err := c.execute(adminMethodSpec{
		name: "OrgDomains",
		graphql: &GraphQLRequest{
			Query:     "query orgDomains($data: ListOrgDomainsRequest!) { _org_domains(params: $data) { " + adminPaginationFields + " org_domains { " + adminOrgDomainFields + " } } }",
			Variables: map[string]interface{}{"data": req},
		},
		graphqlField: "_org_domains",
		restResponse: func() proto.Message { return &authorizerv1.OrgDomainsResponse{} },
		restMethod:   http.MethodPost,
		restPath:     "/v1/admin/org_domains",
		restBody:     req,
		grpcCall: func(ctx context.Context, cli authorizerv1.AuthorizerAdminServiceClient) (interface{}, error) {
			return cli.OrgDomains(ctx, &authorizerv1.OrgDomainsRequest{OrgId: req.GetOrgID(), Pagination: req.protoPagination()})
		},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
