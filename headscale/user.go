package headscale

import "context"
var ()

// ListUsers list all exisint users from Headscale controle plane
func (c *Client)ListUsers()([]string, error){
	var users []string

	return users, nil
}

// CreateUser create a new Headscale user and return true if created by control plane
func (c* Client)CreateUser(ctx context.Context, name string)(bool, error){
	var created bool = false
	queryParam := map[string]string{
		"user": name,
	}
	resp, err := c.get(ctx,"/",queryParam)
	if err != nil {
		return created, err
	}
	defer closeResponseBody(resp)
	created = true
	return created, nil
}

// DeleteUser delete a headscale user from the control plance 
func (c *Client)DeleteUser(ctx context.Context, name string)(bool, error){
	var deleted bool = false
	resp, err := c.delete(ctx, "user/" + name )
	if err != nil {
		return deleted, err
	}
	defer closeResponseBody(resp)
	deleted = true
	return deleted, nil
}