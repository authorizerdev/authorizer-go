package authorizer

import (
	authorizerv1 "github.com/authorizerdev/authorizer-proto-go/authorizer/v1"
)

// Nil-safe accessors for the organization / org SSO / SCIM / org domain request
// types. The gRPC path builds a proto message field by field, and these keep
// that construction safe when the caller passes a nil request (valid for the
// list operations, whose GraphQL params are optional).
func (r *CreateOrganizationRequest) GetName() string {
	if r == nil {
		return ""
	}
	return r.Name
}
func (r *CreateOrganizationRequest) GetDisplayName() *string {
	if r == nil {
		return nil
	}
	return r.DisplayName
}
func (r *UpdateOrganizationRequest) GetID() string {
	if r == nil {
		return ""
	}
	return r.ID
}
func (r *UpdateOrganizationRequest) GetName() *string {
	if r == nil {
		return nil
	}
	return r.Name
}
func (r *UpdateOrganizationRequest) GetDisplayName() *string {
	if r == nil {
		return nil
	}
	return r.DisplayName
}
func (r *UpdateOrganizationRequest) GetEnabled() *bool {
	if r == nil {
		return nil
	}
	return r.Enabled
}
func (r *OrganizationRequest) GetID() string {
	if r == nil {
		return ""
	}
	return r.ID
}
func (r *AddOrgMemberRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}
func (r *AddOrgMemberRequest) GetUserID() string {
	if r == nil {
		return ""
	}
	return r.UserID
}
func (r *AddOrgMemberRequest) GetRoles() []string {
	if r == nil {
		return nil
	}
	return r.Roles
}
func (r *RemoveOrgMemberRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}
func (r *RemoveOrgMemberRequest) GetUserID() string {
	if r == nil {
		return ""
	}
	return r.UserID
}
func (r *ListOrgMembersRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}
func (r *CreateOrgOIDCConnectionRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}
func (r *CreateOrgOIDCConnectionRequest) GetName() string {
	if r == nil {
		return ""
	}
	return r.Name
}
func (r *CreateOrgOIDCConnectionRequest) GetIssuerURL() string {
	if r == nil {
		return ""
	}
	return r.IssuerURL
}
func (r *CreateOrgOIDCConnectionRequest) GetClientID() string {
	if r == nil {
		return ""
	}
	return r.ClientID
}
func (r *CreateOrgOIDCConnectionRequest) GetClientSecret() string {
	if r == nil {
		return ""
	}
	return r.ClientSecret
}
func (r *CreateOrgOIDCConnectionRequest) GetScopes() *string {
	if r == nil {
		return nil
	}
	return r.Scopes
}
func (r *CreateOrgOIDCConnectionRequest) GetRedirectURI() *string {
	if r == nil {
		return nil
	}
	return r.RedirectURI
}
func (r *UpdateOrgOIDCConnectionRequest) GetID() string {
	if r == nil {
		return ""
	}
	return r.ID
}
func (r *UpdateOrgOIDCConnectionRequest) GetName() *string {
	if r == nil {
		return nil
	}
	return r.Name
}
func (r *UpdateOrgOIDCConnectionRequest) GetIssuerURL() *string {
	if r == nil {
		return nil
	}
	return r.IssuerURL
}
func (r *UpdateOrgOIDCConnectionRequest) GetClientID() *string {
	if r == nil {
		return nil
	}
	return r.ClientID
}
func (r *UpdateOrgOIDCConnectionRequest) GetClientSecret() *string {
	if r == nil {
		return nil
	}
	return r.ClientSecret
}
func (r *UpdateOrgOIDCConnectionRequest) GetScopes() *string {
	if r == nil {
		return nil
	}
	return r.Scopes
}
func (r *UpdateOrgOIDCConnectionRequest) GetRedirectURI() *string {
	if r == nil {
		return nil
	}
	return r.RedirectURI
}
func (r *UpdateOrgOIDCConnectionRequest) GetIsActive() *bool {
	if r == nil {
		return nil
	}
	return r.IsActive
}
func (r *OrgOIDCConnectionRequest) GetID() *string {
	if r == nil {
		return nil
	}
	return r.ID
}
func (r *OrgOIDCConnectionRequest) GetOrgID() *string {
	if r == nil {
		return nil
	}
	return r.OrgID
}
func (r *CreateOrgSAMLConnectionRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}
func (r *CreateOrgSAMLConnectionRequest) GetName() string {
	if r == nil {
		return ""
	}
	return r.Name
}
func (r *CreateOrgSAMLConnectionRequest) GetIdpEntityID() string {
	if r == nil {
		return ""
	}
	return r.IdpEntityID
}
func (r *CreateOrgSAMLConnectionRequest) GetIdpSSOURL() string {
	if r == nil {
		return ""
	}
	return r.IdpSSOURL
}
func (r *CreateOrgSAMLConnectionRequest) GetIdpCertificate() string {
	if r == nil {
		return ""
	}
	return r.IdpCertificate
}
func (r *CreateOrgSAMLConnectionRequest) GetSpEntityID() *string {
	if r == nil {
		return nil
	}
	return r.SpEntityID
}
func (r *CreateOrgSAMLConnectionRequest) GetAcsURL() *string {
	if r == nil {
		return nil
	}
	return r.AcsURL
}
func (r *CreateOrgSAMLConnectionRequest) GetAttributeMapping() *string {
	if r == nil {
		return nil
	}
	return r.AttributeMapping
}
func (r *CreateOrgSAMLConnectionRequest) GetAllowIdpInitiated() *bool {
	if r == nil {
		return nil
	}
	return r.AllowIdpInitiated
}
func (r *UpdateOrgSAMLConnectionRequest) GetID() string {
	if r == nil {
		return ""
	}
	return r.ID
}
func (r *UpdateOrgSAMLConnectionRequest) GetName() *string {
	if r == nil {
		return nil
	}
	return r.Name
}
func (r *UpdateOrgSAMLConnectionRequest) GetIdpEntityID() *string {
	if r == nil {
		return nil
	}
	return r.IdpEntityID
}
func (r *UpdateOrgSAMLConnectionRequest) GetIdpSSOURL() *string {
	if r == nil {
		return nil
	}
	return r.IdpSSOURL
}
func (r *UpdateOrgSAMLConnectionRequest) GetIdpCertificate() *string {
	if r == nil {
		return nil
	}
	return r.IdpCertificate
}
func (r *UpdateOrgSAMLConnectionRequest) GetSpEntityID() *string {
	if r == nil {
		return nil
	}
	return r.SpEntityID
}
func (r *UpdateOrgSAMLConnectionRequest) GetAcsURL() *string {
	if r == nil {
		return nil
	}
	return r.AcsURL
}
func (r *UpdateOrgSAMLConnectionRequest) GetAttributeMapping() *string {
	if r == nil {
		return nil
	}
	return r.AttributeMapping
}
func (r *UpdateOrgSAMLConnectionRequest) GetAllowIdpInitiated() *bool {
	if r == nil {
		return nil
	}
	return r.AllowIdpInitiated
}
func (r *UpdateOrgSAMLConnectionRequest) GetIsActive() *bool {
	if r == nil {
		return nil
	}
	return r.IsActive
}
func (r *OrgSAMLConnectionRequest) GetID() *string {
	if r == nil {
		return nil
	}
	return r.ID
}
func (r *OrgSAMLConnectionRequest) GetOrgID() *string {
	if r == nil {
		return nil
	}
	return r.OrgID
}
func (r *CreateScimEndpointRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}
func (r *ScimEndpointRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}
func (r *UserOrganizationsRequest) GetUserID() string {
	if r == nil {
		return ""
	}
	return r.UserID
}
func (r *RequestOrgDomainRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}
func (r *RequestOrgDomainRequest) GetDomain() string {
	if r == nil {
		return ""
	}
	return r.Domain
}
func (r *VerifyOrgDomainRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}
func (r *VerifyOrgDomainRequest) GetDomain() string {
	if r == nil {
		return ""
	}
	return r.Domain
}
func (r *AddVerifiedOrgDomainRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}
func (r *AddVerifiedOrgDomainRequest) GetDomain() string {
	if r == nil {
		return ""
	}
	return r.Domain
}
func (r *DeleteOrgDomainRequest) GetDomain() string {
	if r == nil {
		return ""
	}
	return r.Domain
}
func (r *ListOrgDomainsRequest) GetOrgID() string {
	if r == nil {
		return ""
	}
	return r.OrgID
}

// protoPagination converts the SDK pagination input to its proto equivalent,
// tolerating both a nil request and a nil pagination field.
func (r *ListOrganizationsRequest) protoPagination() *authorizerv1.PaginationRequest {
	if r == nil || r.Pagination == nil {
		return nil
	}
	return &authorizerv1.PaginationRequest{Page: r.Pagination.Page, Limit: r.Pagination.Limit}
}
func (r *ListOrgMembersRequest) protoPagination() *authorizerv1.PaginationRequest {
	if r == nil || r.Pagination == nil {
		return nil
	}
	return &authorizerv1.PaginationRequest{Page: r.Pagination.Page, Limit: r.Pagination.Limit}
}
func (r *UserOrganizationsRequest) protoPagination() *authorizerv1.PaginationRequest {
	if r == nil || r.Pagination == nil {
		return nil
	}
	return &authorizerv1.PaginationRequest{Page: r.Pagination.Page, Limit: r.Pagination.Limit}
}
func (r *ListOrgDomainsRequest) protoPagination() *authorizerv1.PaginationRequest {
	if r == nil || r.Pagination == nil {
		return nil
	}
	return &authorizerv1.PaginationRequest{Page: r.Pagination.Page, Limit: r.Pagination.Limit}
}
