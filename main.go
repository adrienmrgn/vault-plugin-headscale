package main

import (
	"os"

	"github.com/adrienmrgn/vault-plugin-headscale/backend"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
	
)

var (
	exitCode = 1
)

func run(logger hclog.Logger) error {
	meta := &api.PluginAPIClientMeta{}

	err := meta.FlagSet().Parse(os.Args[1:])

	if err != nil {
		return err
	}

	tlsConfig := meta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)
	err = plugin.Serve(&plugin.ServeOpts{
		TLSProviderFunc: 		tlsProviderFunc,
		Logger: 						logger,
		BackendFactoryFunc: backend.Factory,
	})
	if err != nil {
		return err
	}
	return nil 
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{})

	err := run(logger)

	if (err != nil) {
		logger.Error("Error initialising plugin headscale", "error", err)
		os.Exit(exitCode)
	}
}

