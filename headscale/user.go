package headscale

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

// UserStatus : defined the status of a Headscale user
type UserStatus uint8

// Defined UserStatus enum values
const (
	UserCreated UserStatus = iota
	UserExists  UserStatus = iota
	UserDeleted UserStatus = iota
	UserUnknown UserStatus = iota
	UserError   UserStatus = iota
)

// UserConfig : struct that defines a Headscale users
type UserConfig struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// ListUsers list all exisint users from Headscale controle plane
func (c *Client) ListUsers(ctx context.Context) (users []UserConfig, err error) {
	resp, err := c.get(ctx, "/user", nil)
	defer closeResponseBody(resp)
	if err != nil {
		return users, err
	}
	return checkUsersList(resp)
}

func checkUsersList(response *http.Response) (users []UserConfig, err error) {

	switch response.StatusCode {
	case http.StatusOK:
		respBody, err := io.ReadAll(response.Body)
		if err != nil {
			return []UserConfig{}, err
		}
		err = json.Unmarshal(respBody, &users)
		if err != nil {
			return []UserConfig{}, err
		}
		return users, nil
	}
	return []UserConfig{}, err
}

// GetUser : return a Headscale user and its status
func (c *Client) GetUser(ctx context.Context, name string) (status UserStatus, user UserConfig, err error) {
	resp, err := c.get(ctx, "/user/"+name, nil)
	return checkUserGetStatus(resp)
}

func checkUserGetStatus(response *http.Response) (status UserStatus, user UserConfig, err error) {

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return UserError, UserConfig{}, err
	}
	switch response.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &user)
		return UserExists, user, nil
	case http.StatusInternalServerError:
		isMessageUserAlreadyExists := strings.Contains(string(body), "User already exists")
		if isMessageUserAlreadyExists {
			return UserExists, UserConfig{}, nil
		}
	}
	return UserError, UserConfig{}, nil
}

// CreateUser create a new Headscale user and return its status
func (c *Client) CreateUser(ctx context.Context, name string) (status UserStatus, user UserConfig, err error) {

	var requestBody = make(map[string]string)
	requestBody["name"] = name
	resp, err := c.post(ctx, "/user", requestBody)
	defer closeResponseBody(resp)
	if err != nil {
		return UserError, UserConfig{}, err
	}
	return checkUserCreationStatus(resp)
}

func checkUserCreationStatus(response *http.Response) (UserStatus, UserConfig, error) {

	var user UserConfig
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return UserError, UserConfig{}, err
	}
	switch response.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &user)
		return UserCreated, user, nil
	case http.StatusInternalServerError:
		isMessageUserAlreadyExists := strings.Contains(string(body), "User already exists")
		if isMessageUserAlreadyExists {
			return UserExists, UserConfig{}, nil
		}
	}
	return UserError, UserConfig{}, nil
}

// DeleteUser delete a headscale user from the control plance and return deletion status
func (c *Client) DeleteUser(ctx context.Context, name string) (status UserStatus, err error) {
	status = UserUnknown
	resp, err := c.delete(ctx, "/user/"+name)
	defer closeResponseBody(resp)
	if err != nil {
		return UserError, err
	}
	defer closeResponseBody(resp)
	return checkUserDeletionStatus(resp)
}

func checkUserDeletionStatus(response *http.Response) (status UserStatus, err error) {
	switch response.StatusCode {
	case http.StatusOK:
		return UserDeleted, nil
	case http.StatusInternalServerError:
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return UserError, err
		}
		isMessageUserNotFound := strings.Contains(string(body), "User not found")
		if isMessageUserNotFound {
			return UserUnknown, ErrUserNotFound
		}
	}
	return UserError, nil
}
