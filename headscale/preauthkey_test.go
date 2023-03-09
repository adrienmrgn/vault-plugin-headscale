package headscale

import (
	"context"
	"time"

	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePreAuthKey(t *testing.T) {

	c := NewClient()
	c.APIURL = "http://localhost:8080"
	c.APIKey = "sXqH2YEY7Q.Lj-d4ywYsGWLhCiYk9oFIYM-ZrSrq07uxBMpaTbrB_s"
	existingUserName := "bar"
	nonExistingUserName := "baz"
	c.CreateUser(context.Background(), existingUserName)
	c.DeleteUser(context.Background(), nonExistingUserName)
	testCases := []struct {
		pakConfig            PreAuthKeyConfig
		name                 string
		wantError            error
		wantPreAuthKeyStatus PreAuthKeuStatus
	}{
		{
			name: "Simplest request",
			pakConfig: PreAuthKeyConfig{
				User: existingUserName,
			},
			wantError:            nil,
			wantPreAuthKeyStatus: preAuthKeyCreated,
		},
		{
			name: "User does not exists",
			pakConfig: PreAuthKeyConfig{
				User: nonExistingUserName,
			},
			wantError:            ErrUserNotFound,
			wantPreAuthKeyStatus: preAuthKeyError,
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
			wantPreAuthKeyStatus: preAuthKeyCreated,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			preAuthKeyStatus, _, err := c.CreatePreAuthKey(context.Background(), tc.pakConfig)
			assert.ErrorIs(t, err, tc.wantError)
			assert.Equal(t, tc.wantPreAuthKeyStatus, preAuthKeyStatus)
		})
	}
}
