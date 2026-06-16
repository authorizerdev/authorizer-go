package authorizer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// outgoingContext attaches the given headers as gRPC outgoing metadata.
func outgoingContext(ctx context.Context, headers map[string]string) context.Context {
	if len(headers) == 0 {
		return ctx
	}
	md := metadata.MD{}
	for k, v := range headers {
		md.Set(k, v)
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// executeREST performs an HTTP request against a typed REST endpoint.
//   - method: http.MethodPost or http.MethodGet.
//   - path: the REST path from the proto annotation, e.g. "/v1/login".
//   - body: request payload (marshalled to JSON for POST; ignored for GET).
//   - extraHeaders + perCallHeaders are merged onto the request.
//   - out: pointer the JSON response body is unmarshalled into (may be nil).
//
// The Origin header is auto-injected for the same CSRF reason as ExecuteGraphQL.
func (c *AuthorizerClient) executeREST(method, path string, body interface{}, perCallHeaders map[string]string, out interface{}) error {
	return doREST(c.AuthorizerURL, method, path, body, c.ExtraHeaders, perCallHeaders, out)
}

// doREST is the shared REST executor used by both the user and admin clients.
func doREST(baseURL, method, path string, body interface{}, extraHeaders, perCallHeaders map[string]string, out interface{}) error {
	var reqBody io.Reader
	if method == http.MethodPost && body != nil {
		jsonReq, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(jsonReq)
	}

	httpReq, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	for key, val := range extraHeaders {
		httpReq.Header.Add(key, val)
	}
	for key, val := range perCallHeaders {
		httpReq.Header.Add(key, val)
	}

	// Authorizer's CSRF guard rejects state-changing requests without an
	// Origin or Referer header (see ExecuteGraphQL for the full rationale).
	if httpReq.Header.Get("Origin") == "" {
		if u, err := url.Parse(baseURL); err == nil && u.Scheme != "" && u.Host != "" {
			httpReq.Header.Set("Origin", u.Scheme+"://"+u.Host)
		}
	}

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		// Typed REST errors come back as {"message": "..."} or a google.rpc
		// status. Surface the message when present, else the raw body.
		var errRes struct {
			Message string `json:"message"`
			Error   string `json:"error"`
		}
		if json.Unmarshal(bodyBytes, &errRes) == nil {
			if errRes.Message != "" {
				return errors.New(errRes.Message)
			}
			if errRes.Error != "" {
				return errors.New(errRes.Error)
			}
		}
		return errors.New(http.StatusText(res.StatusCode) + ": " + string(bodyBytes))
	}

	if out == nil {
		return nil
	}
	// Proto-typed responses must be decoded with protojson: grpc-gateway emits
	// int64/uint64 fields as JSON strings, which encoding/json cannot unmarshal
	// into Go int64 fields.
	if msg, ok := out.(proto.Message); ok {
		return protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(bodyBytes, msg)
	}
	return json.Unmarshal(bodyBytes, out)
}
