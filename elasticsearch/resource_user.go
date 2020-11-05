package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"

	api "github.com/elastic/go-elasticsearch/v7"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreateOrUpdate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserCreateOrUpdate,
		DeleteContext: resourceUserDelete,
		Schema:        userResource.Schema,
	}
}

func resourceUserCreateOrUpdate(context context.Context, data *schema.ResourceData, state interface{}) diag.Diagnostics {
	client := state.(*api.Client)

	user := userModel{
		Username: data.Get(usernameKey).(string),
		Password: data.Get(passwordKey).(string),
		Enabled:  data.Get(enabledKey).(bool),
		Roles:    []string{},
	}

	if email, exists := data.GetOk(emailKey); exists {
		user.Email = email.(string)
	}

	if fullName, exists := data.GetOk(fullNameKey); exists {
		user.FullName = fullName.(string)
	}

	if roles, exists := data.GetOk(rolesKey); exists {
		user.Roles = mapStringArray(roles.([]interface{}))
	}

	if metadata, exists := data.GetOk(metadataKey); exists {
		user.Metadata = mapStringMap(metadata.(map[string]interface{}))
	}

	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(user); err != nil {
		return diag.FromErr(err)
	}

	response, err := client.Security.PutUser(user.Username, &buffer)
	if err != nil {
		return diag.FromErr(err)
	}

	defer response.Body.Close()

	if response.IsError() {
		return diag.Errorf("Failed to create user: [%d] %s", response.StatusCode, response.String())
	}

	// Empty the response body...
	io.Copy(ioutil.Discard, response.Body)

	data.SetId(user.Username)

	return resourceUserRead(context, data, state)
}

func resourceUserRead(context context.Context, data *schema.ResourceData, state interface{}) diag.Diagnostics {
	client := state.(*api.Client)

	var diags diag.Diagnostics

	username := data.Id()

	response, err := client.Security.GetUser(client.Security.GetUser.WithUsername(username))
	if err != nil {
		return diag.FromErr(err)
	}

	defer response.Body.Close()

	var usersResponse map[string]userModel

	err = json.NewDecoder(response.Body).Decode(&usersResponse)

	if err != nil {
		return diag.FromErr(err)
	}

	user, exists := usersResponse[username]

	if exists {
		if err == nil {
			err = data.Set(usernameKey, user.Username)
		}

		if err == nil {
			err = data.Set(enabledKey, user.Enabled)
		}

		if err == nil {
			err = data.Set(fullNameKey, user.FullName)
		}

		if err == nil {
			err = data.Set(rolesKey, user.Roles)
		}

		if err == nil {
			err = data.Set(metadataKey, user.Metadata)
		}
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserDelete(context context.Context, data *schema.ResourceData, state interface{}) diag.Diagnostics {
	client := state.(*api.Client)

	var diags diag.Diagnostics

	username := data.Id()

	response, err := client.Security.DeleteUser(username)

	if err != nil {
		return diag.FromErr(err)
	}

	defer response.Body.Close()

	if response.IsError() {
		return diag.Errorf("Failed to delete user: [%d] %s", response.StatusCode, response.String())
	}

	// Empty the response body
	io.Copy(ioutil.Discard, response.Body)

	return diags
}
