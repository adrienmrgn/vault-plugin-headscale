package backend

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type headscaleConfig struct {
	APIKey string `json:"api_key"`
	APIURL string `json:"api_url"`
}

func (hc *headscaleConfig) toLogical() *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			"api_url": hc.APIURL,
			"api_key": hc.APIKey,
		},
	}
}

func pathConfigAccess(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: configPath,
		Fields: map[string]*framework.FieldSchema{
			"api_url": {
				Type:        framework.TypeString,
				Description: apiURLDescr,
			},
			"api_key": {
				Type:        framework.TypeString,
				Description: apiKeyDescr,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:    b.ReadHeadscaleConfig,
				Description: readConfigDescr,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback:    b.UpdateHeadscaleConfig,
				Description: updateConfigDescr,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.UpdateHeadscaleConfig,
			},
		},
	}
}

func (b *backend) UpdateHeadscaleConfig(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// build configuration for headscale backend
	config := headscaleConfig{
		APIKey: data.Get("api_key").(string),
		APIURL: data.Get("api_url").(string),
	}

	// build storage entry from configuration
	entry, err := logical.StorageEntryJSON(configPath, config)
	if err != nil {
		return nil, err
	}

	// store configuration in the backend
	err = request.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) ReadHeadscaleConfig(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	headscaleConfig, err := b.retrieveHeadscaleConfig(ctx, request)

	switch {
	case err != nil:
		return nil, err
	case headscaleConfig == nil:
		errorResp := fmt.Sprintf("access configuration for Headscale plugin not configured at %s", configPath)
		return logical.ErrorResponse(errorResp), nil
	}

	return headscaleConfig.toLogical(), nil
}

func (b *backend) retrieveHeadscaleConfig(ctx context.Context, request *logical.Request) (*headscaleConfig, error) {
	entry, err := request.Storage.Get(ctx, configPath)
	switch {
	case err != nil:
		return nil, err
	case entry == nil:
		return &headscaleConfig{}, nil
	}
	config := &headscaleConfig{}
	err = entry.DecodeJSON(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
