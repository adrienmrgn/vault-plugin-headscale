package headscale

import (
	"context"
	"time"

	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePreAuthKey(t *testing.T) {

	client, container, err  := runHeadscale()
	defer container.Terminate()
	assert.NoError(t, err)
	
	existingUserName := "bar"
	nonExistingUserName := "baz"
	client.CreateUser(context.Background(), existingUserName)
	client.DeleteUser(context.Background(), nonExistingUserName)
	testCases := []struct {
		pakConfig            PreAuthKeyConfig
		name                 string
		wantError            error
		wantPreAuthKeyStatus PreAuthKeyStatus
	}{
		{
			name: "Simplest request",
			pakConfig: PreAuthKeyConfig{
				User: existingUserName,
			},
			wantError:            nil,
			wantPreAuthKeyStatus: PreAuthKeyCreated,
		},
		{
			name: "User does not exists",
			pakConfig: PreAuthKeyConfig{
				User: nonExistingUserName,
			},
			wantError:            ErrUserNotFound,
			wantPreAuthKeyStatus: PreAuthKeyError,
		},
		{
			name: "Rrquest With all parameters",
			pakConfig: PreAuthKeyConfig{
				User:       existingUserName,
				Reusable:   true,
				Ephemeral:  false,
				Expiration: time.Now().Add(time.Hour),
				Tags:       []string{"hello", "world"},
			},
			wantError:            nil,
			wantPreAuthKeyStatus: PreAuthKeyCreated,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			preAuthKeyStatus, _, err := client.CreatePreAuthKey(context.Background(), tc.pakConfig)
			assert.ErrorIs(t, err, tc.wantError)
			assert.Equal(t, tc.wantPreAuthKeyStatus, preAuthKeyStatus)
		})
	}
}
