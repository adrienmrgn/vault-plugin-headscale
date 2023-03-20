package headscale

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	apiURI      = "/api/v1"
	contentType = "application/json"
)

type httpMethod string

const (
	httpGet    httpMethod = http.MethodGet
	httpPost   httpMethod = http.MethodPost
	httpDelete httpMethod = http.MethodDelete
)

type requestOptions struct {
	context     context.Context
	uri         string
	method      httpMethod
	queryParams map[string]string
	queryBody   any
}

func (c *Client) buildHTTPRequest(reqOpt requestOptions) (*http.Request, error) {

	var err error
	// convert body to json buffer
	var bodyBytes []byte
	if reqOpt.queryBody != nil {
		bodyBytes, err = json.Marshal(reqOpt.queryBody)
	}
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	// Create HTTP request
	request, err := http.NewRequestWithContext(
		reqOpt.context,
		string(reqOpt.method),
		c.APIURL+apiURI+reqOpt.uri,
		bodyBuffer,
	)
	if err != nil {
		return nil, err
	}

	// Add parameteers
	reqParams := url.Values{}
	for key, value := range reqOpt.queryParams {
		reqParams.Add(key, value)
	}
	request.URL.RawQuery = reqParams.Encode()

	// Add headers
	if len(c.APIKey) > 0 {
		request.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	switch {
	case reqOpt.queryBody == nil:
		request.Header.Set("Accept", contentType)
	default:
		request.Header.Set("Content-Type", contentType)
	}

	return request, nil
}

func (c *Client) get(ctx context.Context, uri string, queryParams map[string]string) (*http.Response, error) {
	reqOpt := requestOptions{
		context:     ctx,
		uri:         uri,
		method:      httpGet,
		queryParams: queryParams,
		queryBody:   nil,
	}
	request, err := c.buildHTTPRequest(reqOpt)
	if err != nil {
		return nil, err
	}
	// Send request
	resp, err := c.HTTP.Do(request)
	if err != nil {
		closeResponseBody(resp)
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	return resp, nil
}

func (c *Client) post(ctx context.Context, uri string, queryBody any) (*http.Response, error) {
	reqOpt := requestOptions{
		context:     ctx,
		uri:         uri,
		method:      httpPost,
		queryParams: nil,
		queryBody:   queryBody,
	}
	request, err := c.buildHTTPRequest(reqOpt)
	if err != nil {
		return nil, err
	}
	// Send request
	resp, err := c.HTTP.Do(request)
	if err != nil {
		closeResponseBody(resp)
		return nil, err
	}

	return resp, nil
}

func (c *Client) delete(ctx context.Context, uri string) (*http.Response, error) {
	reqOpt := requestOptions{
		context:     ctx,
		uri:         uri,
		method:      httpDelete,
		queryParams: nil,
		queryBody:   nil,
	}
	request, err := c.buildHTTPRequest(reqOpt)
	if err != nil {
		return nil, err
	}
	// Send request
	resp, err := c.HTTP.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// instpired by cHarshicorp's consul golang api
func closeResponseBody(resp *http.Response) error {
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.Body.Close()
}
