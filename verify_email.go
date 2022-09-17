package authorizer

import (
	"encoding/json"
	"fmt"
)

// VerifyEmailInput defines attributes for verify_email request
type VerifyEmailInput struct {
	Token string `json:"token"`
}

// VerifyEmail is method attached to AuthorizerClient.
// It performs verify_email mutation on authorizer instance.
// It returns AuthTokenResponse reference or error.
// For implementation details check VerifyEmailExample examples/verify_email.go
func (c *AuthorizerClient) VerifyEmail(req *VerifyEmailInput) (*AuthTokenResponse, error) {
	bytesData, err := c.ExecuteGraphQL(&GraphQLRequest{
		Query: fmt.Sprintf(`mutation verifyEmail($data: VerifyEmailInput!) { verify_email(params: $data) { %s}}`, AuthTokenResponseFragment),
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
