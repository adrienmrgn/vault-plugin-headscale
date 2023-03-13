package headscale

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	c := NewClient()
	c.APIURL = "http://localhost:8080"
	c.APIKey = "6mcpYts8WQ.dysKz_tXvRkFFlNV2xwNt462dU5zcFI0LdD7QpS2sxY"
	userName := "gofezavr"
	userStatus, _, err := c.CreateUser(context.Background(), userName)
	assert.NoError(t, err)
	expectedUserStatus := []UserStatus{
		UserCreated,
		UserExists,
	}
	assert.Contains(t, expectedUserStatus, userStatus)

	userStatus, _, err = c.CreateUser(context.Background(), userName)
	assert.NoError(t, err)
	expectedUserStatus = []UserStatus{
		UserExists,
	}
	assert.Contains(t, expectedUserStatus, userStatus)
}

func TestDeleteUser(t *testing.T) {
	c := NewClient()
	c.APIURL = "http://localhost:8080"
	c.APIKey = "6mcpYts8WQ.dysKz_tXvRkFFlNV2xwNt462dU5zcFI0LdD7QpS2sxY"
	userName := "bar"
	userStatus, _, _ := c.CreateUser(context.Background(), userName)
	expectedUserStatus := []UserStatus{
		UserCreated,
		UserExists,
	}
	assert.Contains(t, expectedUserStatus, userStatus)
	userStatus, err := c.DeleteUser(context.Background(), userName)
	assert.Nil(t, err)
	expectedUserStatus = []UserStatus{
		UserDeleted,
		UserUnknown,
	}
	assert.Contains(t, expectedUserStatus, userStatus)

	userStatus, err = c.DeleteUser(context.Background(), userName)
	assert.ErrorIs(t, err, ErrUserNotFound)
	expectedUserStatus = []UserStatus{
		UserUnknown,
	}
	assert.Contains(t, expectedUserStatus, userStatus)
}
