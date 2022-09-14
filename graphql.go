package authorizer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

	fmt.Println(string(jsonReq))

	client := http.Client{}
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
	if headers != nil {
		for key, val := range headers {
			httpReq.Header.Add(key, val)
		}
	}

	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	// Need to close the response stream, once response is read.
	// Hence defer close. It will automatically take care of it.
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
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

	dataBytes, err := json.Marshal(gqlRes.Data)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}
