package authorizer

import (
	"fmt"
	"strings"
)

// authorizerClient defines the attributes required to initiate authorizer client
type authorizerClient struct {
	ClientID      string
	AuthorizerURL string
	RedirectURL   string
	ExtraHeaders  map[string]string
}

// NewAuthorizerClient creates an authorizer client instance and returns reference to it
func NewAuthorizerClient(clientID, authorizerURL, redirectURL string, extraHeaders map[string]string) (*authorizerClient, error) {
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

	return &authorizerClient{
		RedirectURL:   strings.TrimSuffix(redirectURL, "/"),
		AuthorizerURL: strings.TrimSuffix(authorizerURL, "/"),
		ClientID:      clientID,
		ExtraHeaders:  extraHeaders,
	}, nil
}
