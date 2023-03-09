package headscale

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

// PreAuthKeuStatus defines the status of a Headscale preauthkey
type PreAuthKeuStatus uint8

// Instantiate the PreAuthKeyStatus enum
const (
	preAuthKeyCreated PreAuthKeuStatus = iota
	preAuthKeyExists  PreAuthKeuStatus = iota
	preAuthKeyDeleted PreAuthKeuStatus = iota
	preAuthKeyUnknown PreAuthKeuStatus = iota
	preAuthKeyError   PreAuthKeuStatus = iota
)

// PreAuthKeyConfig is used to create a preAuthKey
type PreAuthKeyConfig struct {
	User       string    `json:"name"`
	Reusable   bool      `json:"reusable"`
	Ephemeral  bool      `json:"ephemeral"`
	Expiration time.Time `json:"expiration"`
	Tags       []string  `json:"acl_tags"`
}

// PreAuthKeyResponse stores the HEadscale response
type PreAuthKeyResponse struct {
	PreAuthKey struct {
		User       string    `json:"user"`
		ID         string    `json:"id"`
		Key        string    `json:"key"`
		Reusable   bool      `json:"reusable"`
		Ephemeral  bool      `json:"ephemeral"`
		Used       bool      `json:"used"`
		Expiration time.Time `json:"expiration"`
		CreatedAt  time.Time `json:"createdAt"`
		ACLTags    []string  `json:"aclTags"`
	} `json:"preAuthKey"`
}

// timestampToProtobufTimestamp is used to convert time.Time to a google protobuf timestamp compatible format 
func timestampToProtobufTimestamp(t time.Time) string {
	return t.Format("1992-05-07T:%M:%S.%fZ")
}

// CreatePreAuthKey creates a preAuthKey from a PreAuthKeyConfig
func (c *Client) CreatePreAuthKey(ctx context.Context, preAuthKeyConfig PreAuthKeyConfig) (status PreAuthKeuStatus, preAuthKey PreAuthKeyResponse, err error) {

	preAuthKey = PreAuthKeyResponse{}

	requestBody := buildPreAuthKeyRequestBody(preAuthKeyConfig)

	resp, err := c.post(ctx, "/preauthkey", requestBody)
	defer closeResponseBody(resp)

	if err != nil {
		return preAuthKeyError, PreAuthKeyResponse{}, err
	}

	status, err = checkPreAuthKeyCreationStatus(resp)
	if err != nil {
		return preAuthKeyError, PreAuthKeyResponse{}, err
	}

	switch status {
	case preAuthKeyCreated:
		preAuthKey, err = retrievePreAuthKeyResponse(resp)
	}

	if err != nil {
		return preAuthKeyError, preAuthKey, err
	}
	return status, preAuthKey, nil
}

func buildPreAuthKeyRequestBody(preAuthKeyConfig PreAuthKeyConfig) map[string]any {
	var requestBody map[string]any
	requestBody = make(map[string]any)
	if !preAuthKeyConfig.Expiration.IsZero() {
		requestBody["expiration"] = timestampToProtobufTimestamp(preAuthKeyConfig.Expiration)
	}
	if len(preAuthKeyConfig.Tags) != 0 {
		var formatedTags []string
		formatedTags = make([]string, len(preAuthKeyConfig.Tags))
		for i, tag := range preAuthKeyConfig.Tags {
			formatedTags[i] = "tag:" + strings.ToLower(tag)
		}
		requestBody["acl_tags"] = formatedTags
	}
	requestBody["user"] = preAuthKeyConfig.User
	requestBody["expiration"] = preAuthKeyConfig.Expiration
	requestBody["ephemeral"] = preAuthKeyConfig.Ephemeral

	return requestBody
}

func checkPreAuthKeyCreationStatus(response *http.Response) (status PreAuthKeuStatus, err error) {

	switch response.StatusCode {
	case http.StatusOK:
		return preAuthKeyCreated, nil
	case http.StatusInternalServerError:
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return preAuthKeyError, err
		}
		isMessageUserNotFound := strings.Contains(string(body), "User not found")
		if isMessageUserNotFound {
			return preAuthKeyError, ErrUserNotFound
		}
	}
	return preAuthKeyUnknown, nil
}

func retrievePreAuthKeyResponse(response *http.Response) (preAuthKeyResponse PreAuthKeyResponse, err error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return PreAuthKeyResponse{}, err
	}

	var responseData PreAuthKeyResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return PreAuthKeyResponse{}, err
	}
	return responseData, nil
}
