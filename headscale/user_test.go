package headscale

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	client, container, err  := runHeadscale()
	defer container.Terminate()
	userName := "foo"
	userStatus, _, err := client.CreateUser(container.Context, userName)
	assert.NoError(t, err)
	expectedUserStatus := []UserStatus{
		UserCreated,
	}
	assert.Contains(t, expectedUserStatus, userStatus)

	userStatus, _, err = client.CreateUser(container.Context, userName)
	assert.NoError(t, err)
	expectedUserStatus = []UserStatus{
		UserExists,
	}
	assert.Contains(t, expectedUserStatus, userStatus)
}


func TestDeleteUser(t *testing.T) {
	client, container, err  := runHeadscale()
	defer container.Terminate()
	userName := "bar"
	userStatus, _, _ := client.CreateUser(context.Background(), userName)
	expectedUserStatus := []UserStatus{
		UserCreated,
	}
	assert.Contains(t, expectedUserStatus, userStatus)
	userStatus, err = client.DeleteUser(context.Background(), userName)
	assert.Nil(t, err)
	expectedUserStatus = []UserStatus{
		UserDeleted,
	}
	assert.Contains(t, expectedUserStatus, userStatus)

	userStatus, err = client.DeleteUser(context.Background(), userName)
	assert.ErrorIs(t, err, ErrUserNotFound)
	expectedUserStatus = []UserStatus{
		UserUnknown,
	}
	assert.Contains(t, expectedUserStatus, userStatus)
}
