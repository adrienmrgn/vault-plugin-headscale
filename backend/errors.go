package backend

import (
	"errors"
)

// Defines Global errors
var (
	ErrFailedToCreateHeadscaleUser 	= errors.New("failed to create Headscale user")
	ErrEmptyConfigEntry							= errors.New("failed to retrieve config from vault backend")
	ErrDeleteUser										= errors.New("failed to deleted user from Headscale controle plane")
)