package backend

import (
	"context"
	"time"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type headscaleUserConfig struct {
	UserName 			string 		`json:"user_name"`
	CreatedBy 		string 		`json:"created_by"`
	CreationTime 	time.Time	`json:"creation_time"`
}

func pathListUsers(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: userPath+"/?$",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the user.",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation		: &framework.PathOperation{
				Callback: b.ListHeadscaleUsers,
				Description: listUserDescr,
			},
		},
	}
}

func pathUser(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: userPath + "/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the user.",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation		: &framework.PathOperation{
				Callback: b.ReadHeadscaleUser,
				Description: readUserDescr,
			},
			logical.DeleteOperation : &framework.PathOperation{
				Callback: 		b.DeleteHeadscaleUser,
				Description: 	deleteUserDescr,
			},
			logical.CreateOperation : &framework.PathOperation{
				Callback: 		b.UpdateHeadscaleUser,
			},
			logical.UpdateOperation	 : &framework.PathOperation{
				Callback: 		b.UpdateHeadscaleUser,
				Description: 	updateUserDescr,
			},
		},
	}
}

func (b *backend) ListHeadscaleUsers(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entries, err := request.Storage.List(ctx, userPath+"/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) ReadHeadscaleUser(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	entry, err := request.Storage.Get(ctx, userPath+"/"+name)
	if err != nil {
		return logical.ErrorResponse("failed to read data at %s",userPath+"/"+name), err
	}
	if entry == nil {
		return nil, nil
	}

	var headscaleUserConfigData headscaleUserConfig
	err = entry.DecodeJSON(&headscaleUserConfigData)
	if err != nil {
		return logical.ErrorResponse("failed to decode entry as Headscale User Configuration"), err
	}
	response := &logical.Response{
		Data:	map[string]interface{}{
			"user_name": 			headscaleUserConfigData.UserName,
			"create_by": 			headscaleUserConfigData.CreatedBy,
			"creation_time":	headscaleUserConfigData.CreationTime,
		},
	}
	return response, nil
}

func (b *backend) UpdateHeadscaleUser(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	headscaleConfig, err := b.retrieveHeadscaleConfig(ctx, request)
	switch  {
	case err != nil :
		return nil, err
	case headscaleConfig == nil :
		errorResp := fmt.Sprintf("access configuration for Headscale plugin not configured at %s", configPath)
		return logical.ErrorResponse(errorResp), nil
	}
	name := data.Get("name").(string)
	entry, err := logical.StorageEntryJSON(userPath+"/"+name,headscaleUserConfig{
		UserName: name,
		CreatedBy: "vault",
		CreationTime: time.Now(),
	})
	if err != nil {
		return logical.ErrorResponse("failed to build Headscale user entry"), err
	}
	err = request.Storage.Put(ctx, entry)
	if err != nil {
		return logical.ErrorResponse("failed to store Headscale user config"), err
	}
	return nil, nil
}

func (b *backend) DeleteHeadscaleUser(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	headscaleConfig, err := b.retrieveHeadscaleConfig(ctx, request)
	switch  {
	case err != nil :
		return nil, err
	case headscaleConfig == nil :
		errorResp := fmt.Sprintf("access configuration for Headscale plugin not configured at %s", configPath)
		return logical.ErrorResponse(errorResp), nil
	}
	
	name := data.Get("name").(string)
	err = request.Storage.Delete(ctx, userPath+"/"+name); 
	if err != nil {
		return logical.ErrorResponse("failed to delete entry at %s",userPath+"/"+name), err
	}
	return nil, nil
}