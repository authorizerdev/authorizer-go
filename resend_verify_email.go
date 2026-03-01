package authorizer

import (
	"encoding/json"
)

// ResendVerifyEmailRequest defines attributes for resend_verify_email request
type ResendVerifyEmailRequest struct {
	Email       string  `json:"email"`
	Identifier  *string `json:"identifier,omitempty"`
}

// ResendVerifyEmail is method attached to AuthorizerClient.
// It performs resend_verify_email mutation on authorizer instance.
// It takes ResendVerifyEmailRequest reference as parameter and returns Response reference or error.
func (c *AuthorizerClient) ResendVerifyEmail(req *ResendVerifyEmailRequest) (*Response, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: `mutation resendVerifyEmail($data: ResendVerifyEmailRequest!) { resend_verify_email(params: $data) { message }}`,
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*Response
	json.Unmarshal(bytesData, &res)

	return res["resend_verify_email"], nil
}
