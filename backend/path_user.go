package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/adrienmrgn/vault-plugin-headscale/headscale"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type headscaleUserConfig struct {
	UserName     string    `json:"user_name"`
	UserID       uint32    `json:"user_id"`
	CreatedBy    string    `json:"created_by"`
	CreationTime time.Time `json:"creation_time"`
}

func pathListUsers(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: userPath + "/?$",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the user.",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback:    b.ListHeadscaleUsers,
				Description: listUserDescr,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.UpdateHeadscaleUser,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback:    b.UpdateHeadscaleUser,
				Description: updateUserDescr,
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
			logical.ReadOperation: &framework.PathOperation{
				Callback:    b.ReadHeadscaleUser,
				Description: readUserDescr,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback:    b.DeleteHeadscaleUser,
				Description: deleteUserDescr,
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
	path := userPath + "/" + name
	entry, err := request.Storage.Get(ctx, path)
	if err != nil {
		return logical.ErrorResponse("failed to read data at %s", userPath+"/"+name), err
	}
	if entry == nil {
		return logical.ErrorResponse(fmt.Sprintf("empty entry at %s", path)), ErrEmptyConfigEntry
	}

	var headscaleUserConfigData headscaleUserConfig
	err = entry.DecodeJSON(&headscaleUserConfigData)
	if err != nil {
		errorResponse := fmt.Sprintf("failed to decode entry as Headscale User Configuration at %s", path)
		return logical.ErrorResponse(errorResponse), err
	}
	response := &logical.Response{
		Data: map[string]interface{}{
			"user_name":     headscaleUserConfigData.UserName,
			"user_id":       headscaleUserConfigData.UserID,
			"create_by":     headscaleUserConfigData.CreatedBy,
			"creation_time": headscaleUserConfigData.CreationTime,
		},
	}
	return response, nil
}

func (b *backend) UpdateHeadscaleUser(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	headscaleConfig, err := b.retrieveHeadscaleConfig(ctx, request)
	switch {
	case err != nil:
		return nil, err
	case headscaleConfig == nil:
		errorResp := fmt.Sprintf("access configuration for Headscale plugin not configured at %s", configPath)
		return logical.ErrorResponse(errorResp), ErrEmptyConfigEntry
	}
	name := data.Get("name").(string)
	path := userPath + "/" + name

	headscaleClient.APIURL = headscaleConfig.APIURL
	headscaleClient.APIKey = headscaleConfig.APIKey

	status, user, err := headscaleClient.CreateUser(ctx, name)
	if err != nil {
		errorResp := fmt.Sprintf("accessing Headscale control plane failed")
		return logical.ErrorResponse(errorResp), err
	}

	var entry *logical.StorageEntry
	switch status {
	case headscale.UserCreated:
		entry, err = logical.StorageEntryJSON(path, headscaleUserConfig{
			UserName:     name,
			UserID:       user.ID,
			CreatedBy:    "vault",
			CreationTime: time.Now(),
		})
		if err != nil {
			return logical.ErrorResponse("failed to build Headscale user entry"), err
		}
	case headscale.UserExists:
		entry, err = logical.StorageEntryJSON(path, headscaleUserConfig{
			UserName:     user.Name,
			UserID:       user.ID,
			CreatedBy:    "headscale",
			CreationTime: user.CreatedAt,
		})
		if err != nil {
			return logical.ErrorResponse("failed to build Headscale user entry"), err
		}
	case headscale.UserError:
		return logical.ErrorResponse("failed to create Headscale user"), ErrFailedToCreateHeadscaleUser
	}

	err = request.Storage.Put(ctx, entry)
	if err != nil {
		return logical.ErrorResponse("failed to store Headscale user config"), err
	}
	return nil, nil
}

func (b *backend) DeleteHeadscaleUser(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	headscaleConfig, err := b.retrieveHeadscaleConfig(ctx, request)
	switch {
	case err != nil:
		return nil, err
	case headscaleConfig == nil:
		errorResp := fmt.Sprintf("access configuration for Headscale plugin not configured at %s", configPath)
		return logical.ErrorResponse(errorResp), nil
	}

	name := data.Get("name").(string)
	path := userPath + "/" + name

	headscaleClient.APIURL = headscaleConfig.APIURL
	headscaleClient.APIKey = headscaleConfig.APIKey

	status, err := headscaleClient.DeleteUser(ctx, name)
	if err != nil {
		return logical.ErrorResponse("failed to delete Headscale user %s from control plane", name), err
	}
	switch status {
	case headscale.UserDeleted:
		err = request.Storage.Delete(ctx, path)
		if err != nil {
			return logical.ErrorResponse("failed to delete entry at %s", path), err
		}
		return nil, nil
	case headscale.UserUnknown:
		responseMsg := fmt.Sprintf("failed to delete Headscale user %s from control plane", name)
		// TODO : list users and disply in see_also
		return logical.HelpResponse(responseMsg, nil, nil), err
	case headscale.UserError:
		return logical.ErrorResponse("failed to delete Headscale user %s from control plane", name), ErrDeleteUser
	}
	return nil, nil
}
