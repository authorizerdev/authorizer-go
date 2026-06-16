package authorizer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

// GraphQLRequest is object used to make graphql queries
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLError struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}
type GraphQLResponse struct {
	Errors []*GraphQLError `json:"errors"`
	Data   interface{}     `json:"data"`
}

func (c *AuthorizerClient) ExecuteGraphQL(req *GraphQLRequest, headers map[string]string) ([]byte, error) {
	// Marshal it into JSON prior to requesting
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: 30 * time.Second}
	httpReq, err := http.NewRequest(http.MethodPost, c.AuthorizerURL+"/graphql", bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}

	// set the content type for http request
	httpReq.Header.Set("Content-Type", "application/json")

	// set the default extra headers
	for key, val := range c.ExtraHeaders {
		httpReq.Header.Add(key, val)
	}

	// set the headers for this request
	for key, val := range headers {
		httpReq.Header.Add(key, val)
	}

	// Authorizer's CSRF guard rejects state-changing requests without an
	// Origin or Referer header. Browsers send Origin automatically; this
	// server-side client must set it explicitly. The server's own origin
	// always passes the guard's same-origin rule (the default when
	// ALLOWED_ORIGINS is the wildcard). Deployments with an explicit
	// allowlist can override it via ExtraHeaders or per-call headers.
	if httpReq.Header.Get("Origin") == "" {
		if u, err := url.Parse(c.AuthorizerURL); err == nil && u.Scheme != "" && u.Host != "" {
			httpReq.Header.Set("Origin", u.Scheme+"://"+u.Host)
		}
	}

	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	// Need to close the response stream, once response is read.
	// Hence defer close. It will automatically take care of it.
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var gqlRes *GraphQLResponse
	err = json.Unmarshal(bodyBytes, &gqlRes)
	if err != nil {
		return nil, err
	}

	if len(gqlRes.Errors) > 0 {
		return nil, errors.New(gqlRes.Errors[0].Message)
	}

	// Non-GraphQL failures (e.g. the CSRF guard's 403, a proxy error page)
	// carry no "errors" array. Without this check they would surface as a
	// nil result with a nil error and panic the caller.
	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(http.StatusText(res.StatusCode) + ": " + string(bodyBytes))
	}

	dataBytes, err := json.Marshal(gqlRes.Data)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}
