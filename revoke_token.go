package authorizer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// RevokeTokenInput defines attributes for /oauth/revoke request
type RevokeTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}

// RevokeToken is method attached to AuthorizerClient.
// It performs /oauth/revoke api call on authorizer instance.
// It takes RevokeTokenInput reference as parameter and returns Response reference or error.
// For implementation details check RevokeTokenExample examples/revoke_token.go
func (c *AuthorizerClient) RevokeToken(req *RevokeTokenInput) (*Response, error) {
	if req.RefreshToken == "" {
		return nil, errors.New("refresh_token is required")
	}
	// Marshal it into JSON prior to requesting
	jsonReq, err := json.Marshal(map[string]string{
		"refresh_token": req.RefreshToken,
		"client_id":     c.ClientID,
	})
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, c.AuthorizerURL+"/oauth/revoke", bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}

	// set the default extra headers
	for key, val := range c.ExtraHeaders {
		httpReq.Header.Add(key, val)
	}

	httpRes, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	// Need to close the response stream, once response is read.
	// Hence defer close. It will automatically take care of it.
	defer httpRes.Body.Close()

	bodyBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var res *Response
	json.Unmarshal(bodyBytes, &res)

	return res, nil
}
