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

func resourceAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIKeyCreate,
		ReadContext:   resourceAPIKeyRead,
		DeleteContext: resourceAPIKeyDelete,
		Schema:        apiKeyResource.Schema,
	}
}

func resourceAPIKeyCreate(context context.Context, data *schema.ResourceData, state interface{}) diag.Diagnostics {
	client := state.(*api.Client)

	model := apiKeyModel{
		Name: data.Get(nameKey).(string),
	}

	if expiration, ok := data.GetOk(expirationKey); ok {
		model.Expiration = expiration.(string)
	}

	if rolesSource, ok := data.GetOk(roleDescriptorsKey); ok {
		model.RoleDescriptor = map[string]roleModel{}
		for _, source := range rolesSource.([]interface{}) {
			name, role := mapRole(source)
			model.RoleDescriptor[name] = role
		}
	}

	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(model); err != nil {
		return diag.FromErr(err)
	}

	response, err := client.Security.CreateAPIKey(&buffer)
	if err != nil {
		return diag.FromErr(err)
	}

	defer response.Body.Close()

	if response.IsError() {
		return diag.Errorf("Failed to create API key: [%d] %s", response.StatusCode, response.String())
	}

	var apiKey apiKeyCreateResponse
	json.NewDecoder(response.Body).Decode(&apiKey)

	data.SetId(apiKey.ID)
	data.Set(apiKeyKey, apiKey.APIKey)
	data.Set(nameKey, apiKey.Name)

	var diags diag.Diagnostics

	return diags
}

func resourceAPIKeyRead(context context.Context, data *schema.ResourceData, state interface{}) diag.Diagnostics {
	client := state.(*api.Client)

	var diags diag.Diagnostics

	apiKeyID := data.Id()

	response, err := client.Security.GetAPIKey(client.Security.GetAPIKey.WithID(apiKeyID))

	if err != nil {
		return diag.FromErr(err)
	}

	defer response.Body.Close()

	if response.IsError() {
		return diag.Errorf("Failed to read role: [%d] %s", response.StatusCode, response.String())
	}

	var getReponse apiKeyGetResponse
	if err = json.NewDecoder(response.Body).Decode(&getReponse); err != nil {
		return diag.FromErr(err)
	}

	if len(getReponse.APIKeys) != 1 {
		return diag.Errorf("API key not found")
	}

	apiKey := getReponse.APIKeys[0]

	data.SetId(apiKey.ID)
	data.Set(apiKeyKey, apiKey.APIKey)
	data.Set(nameKey, apiKey.Name)

	return diags
}

func resourceAPIKeyDelete(context context.Context, data *schema.ResourceData, state interface{}) diag.Diagnostics {
	client := state.(*api.Client)

	var diags diag.Diagnostics

	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(map[string]string{
		"id": data.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := client.Security.InvalidateAPIKey(&buffer)

	if err != nil {
		return diag.FromErr(err)
	}

	defer response.Body.Close()

	if response.IsError() {
		return diag.Errorf("Failed to invalidate API key: [%d] %s", response.StatusCode, response.String())
	}

	io.Copy(ioutil.Discard, response.Body)

	return diags
}

func mapRole(item interface{}) (string, roleModel) {
	roleSource := item.(map[string]interface{})

	name := roleSource[nameKey].(string)

	role := roleModel{
		Cluster:      mapStringArray(roleSource[clusterKey].([]interface{})),
		Applications: mapApplications(roleSource[applicationsKey].([]interface{})),
		Indices:      mapIndices(roleSource[indicesKey].([]interface{})),
	}
	return name, role
}
