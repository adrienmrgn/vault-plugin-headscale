package headscale

import (
	"errors"
)

var (
	// ErrUserNotFound : error returned when a headscale user is not found
	ErrUserNotFound = errors.New("User not found")
)