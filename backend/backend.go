package backend

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var Version = ""

type backend struct {
	*framework.Backend
}

const (
	configPath					= "config/access"
	userPath						= "user"
	backendHelp					= "The Headscale backend serves preauthkeys"
	apiKeyDescr 				= "API key used to authenticate backedn to the HEadscale Controle Plane"
	apiURLDescr 				= "API URL of the Headscale Controle Plance"
	updateConfigDescr 	= "Update the Headscale access configuration"
	readConfigDescr			= "Read the Headscale access configuration"
	listUserDescr				= "List headscale users configured from Vault"
	readUserDescr				= "Read a headscale user"
	updateUserDescr			=	"Update a headscale user"
	deleteUserDescr			= "Delete a headscale user"
)

func  Factory(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	err := b.Setup(ctx, config) 
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	b := &backend{}
	b.Backend = &framework.Backend{
		BackendType: 		logical.TypeLogical,
		Help:						strings.TrimSpace(backendHelp),
		RunningVersion: Version,
		PathsSpecial: 	&logical.Paths{
			SealWrapStorage: []string{
				configPath,
			},
		},
		Paths: 					[]*framework.Path{
			pathConfigAccess(b),
			pathListUsers(b),
			pathUser(b),
		},
	}
	return b
}