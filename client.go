// Package authorizer provides client and methods using which you can
// perform various graphql operations to your authorizer instance
package authorizer

import (
	"fmt"
	"strings"
)

// AuthorizerClient defines the attributes required to initiate authorizer client
type AuthorizerClient struct {
	ClientID      string
	AuthorizerURL string
	RedirectURL   string
	ExtraHeaders  map[string]string
}

// NewAuthorizerClient creates an authorizer client instance.
// It returns reference to authorizer client instance or error.
func NewAuthorizerClient(clientID, authorizerURL, redirectURL string, extraHeaders map[string]string) (*AuthorizerClient, error) {
	if strings.TrimSpace(clientID) == "" {
		return nil, fmt.Errorf("clientID missing")
	}

	if strings.TrimSpace(authorizerURL) == "" {
		return nil, fmt.Errorf("authorizerURL missing")
	}

	// extraHeaders is optional parameter,
	// hence if not set, initialize it with empty map
	headers := extraHeaders
	if headers == nil {
		headers = make(map[string]string)
	}

	// if x-authorizer-url is not present
	// set it to authorizerURL
	if _, ok := headers["x-authorizer-url"]; !ok {
		headers["x-authorizer-url"] = authorizerURL
	}

	return &AuthorizerClient{
		RedirectURL:   strings.TrimSuffix(redirectURL, "/"),
		AuthorizerURL: strings.TrimSuffix(authorizerURL, "/"),
		ClientID:      clientID,
		ExtraHeaders:  extraHeaders,
	}, nil
}
