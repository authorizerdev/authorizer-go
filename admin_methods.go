package authorizer

import (
	"context"
	"net/http"

	authorizerv1 "github.com/authorizerdev/authorizer-go/internal/genpb/authorizer/v1"
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
		name:       "AdminLogout",
		restMethod: http.MethodPost,
		restPath:   "/v1/admin/logout",
		restBody:   &authorizerv1.AdminLogoutRequest{},
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
		name:       "AdminSession",
		restMethod: http.MethodGet,
		restPath:   "/v1/admin/session",
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
		name:       "AdminMeta",
		restMethod: http.MethodGet,
		restPath:   "/v1/admin/meta",
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
			Query:     "query users($data: PaginatedRequest) { _users(params: $data) { " + adminPaginationFields + " users { " + adminUserFields + " } } }",
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
			Query:     "query verificationRequests($data: PaginatedRequest) { _verification_requests(params: $data) { " + adminPaginationFields + " verification_requests { " + adminVerifReqFields + " } } }",
			Variables: map[string]interface{}{"data": req},
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
			Query:     "query webhooks($data: PaginatedRequest) { _webhooks(params: $data) { " + adminPaginationFields + " webhooks { " + adminWebhookFields + " } } }",
			Variables: map[string]interface{}{"data": req},
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
			Query:     "query emailTemplates($data: PaginatedRequest) { _email_templates(params: $data) { " + adminPaginationFields + " email_templates { " + adminEmailTplFields + " } } }",
			Variables: map[string]interface{}{"data": req},
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
		name:       "FgaGetModel",
		restMethod: http.MethodGet,
		restPath:   "/v1/admin/fga/model",
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
		name:       "FgaReset",
		restMethod: http.MethodPost,
		restPath:   "/v1/admin/fga/reset",
		restBody:   &authorizerv1.FgaResetRequest{},
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
