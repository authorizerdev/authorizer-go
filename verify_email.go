package authorizer

import (
	"encoding/json"
	"fmt"
)

// VerifyEmailRequest defines attributes for verify_email request
type VerifyEmailRequest struct {
	Token string  `json:"token"`
	State *string `json:"state,omitempty"`
}

// VerifyEmailInput is deprecated: Use VerifyEmailRequest instead
type VerifyEmailInput = VerifyEmailRequest

// VerifyEmail is method attached to AuthorizerClient.
// It performs verify_email mutation on authorizer instance.
// It returns AuthTokenResponse reference or error.
func (c *AuthorizerClient) VerifyEmail(req *VerifyEmailRequest) (*AuthTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`mutation verifyEmail($data: VerifyEmailRequest!) { verify_email(params: $data) { %s}}`, AuthTokenResponseFragment),
		Variables: map[string]interface{}{
			"data": req,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]*AuthTokenResponse
	json.Unmarshal(bytesData, &res)

	return res["verify_email"], nil
}
