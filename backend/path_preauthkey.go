package backend

import (
	"context"
	"fmt"
	"time"

	headscale "github.com/adrienmrgn/headscale-client/client"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	preAuthKeyPath = "creds"
)
func pathPreAuthKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: preAuthKeyPath+ "/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the user.",
			},
			"ephemeral": {
				Type:        framework.TypeBool,
				Description: "Ephemeral preAuthKey",
			},
			"reusable": {
				Type:        framework.TypeBool,
				Description: "Reusable preAuthKey",
			},
			"expiration": {
				Type:        framework.TypeTime,
				Description: "Expiration date",
			},
			"tags": {
				Type:        framework.TypeCommaStringSlice,
				Description: "List of tags applied on preAuthKey",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:    b.CreateHeadscalPreAuthKey,
				Description: CreatePreAuthKey,
			},
		},
	}
}

func (b *backend)CreateHeadscalPreAuthKey(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error){
	
	pakConfig := generatePreAuthKeyConfig(data)

	headscaleConfig, err := b.retrieveHeadscaleConfig(ctx, request)
	switch {
	case err != nil:
		return nil, err
	case headscaleConfig == nil:
		errorResp := fmt.Sprintf("access configuration for Headscale plugin not configured at %s", configPath)
		return logical.ErrorResponse(errorResp), ErrEmptyConfigEntry
	}

	headscaleClient.APIURL = headscaleConfig.APIURL
	headscaleClient.APIKey = headscaleConfig.APIKey

	status, preAuthKey, err := headscaleClient.CreatePreAuthKey(ctx, pakConfig)
	if err != nil {
		errorResp := fmt.Sprintf("accessing Headscale control plane failed")
		return logical.ErrorResponse(errorResp), err
	}
	var response *logical.Response
	switch status {
	case headscale.PreAuthKeyCreated:
		response = &logical.Response{
		Data: map[string]interface{}{
			"user":     		preAuthKey.PreAuthKey.User,
			"id":       		preAuthKey.PreAuthKey.ID,
			"created_at":   preAuthKey.PreAuthKey.CreatedAt,
			"ephemeral": 		preAuthKey.PreAuthKey.Ephemeral,
			"reusable": 		preAuthKey.PreAuthKey.Reusable,
			"expiration": 	preAuthKey.PreAuthKey.Expiration,
			"tags": 				pakConfig.Tags,
			"key": 					preAuthKey.PreAuthKey.Key,
		},
	}
	default:
		response.AddWarning("unhandled case")
}
return response, nil

}

func generatePreAuthKeyConfig(data *framework.FieldData) headscale.PreAuthKeyConfig{
	
	var pakConfig headscale.PreAuthKeyConfig
	name, ok := data.GetOk("name")
	if ok {
		pakConfig.User = name.(string)
	}
	reusable, ok := data.GetOk("reusable")
	if ok {
		pakConfig.Reusable = reusable.(bool)
	}
	ephermeral, ok := data.GetOk("ephemeral")
	if ok {
		pakConfig.Ephemeral = ephermeral.(bool)
	}

	tags, ok := data.GetOk("tags")
	if ok {
		pakConfig.Tags = tags.([]string)
	}
	expiration, ok := data.GetOk("expiration")
	if ok && ! expiration.(time.Time).IsZero() {
		pakConfig.Expiration = expiration.(time.Time)
	}
	return pakConfig
}