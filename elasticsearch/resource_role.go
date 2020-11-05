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

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreateOrUpdate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleCreateOrUpdate,
		DeleteContext: resourceRoleDelete,
		Schema:        roleResource.Schema,
	}
}

func resourceRoleCreateOrUpdate(context context.Context, data *schema.ResourceData, state interface{}) diag.Diagnostics {
	client := state.(*api.Client)

	roleName := data.Get(nameKey).(string)
	role := roleModel{
		Cluster:      mapStringArray(data.Get(clusterKey).([]interface{})),
		Applications: mapApplications(data.Get(applicationsKey).([]interface{})),
		Indices:      mapIndices(data.Get(indicesKey).([]interface{})),
	}

	if metadata, exists := data.GetOk(metadataKey); exists {
		role.Metadata = mapStringMap(metadata.(map[string]interface{}))
	}

	if runAs, exists := data.GetOk(runAsKey); exists {
		role.RunAs = mapStringArray(runAs.([]interface{}))
	}

	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(role); err != nil {
		return diag.FromErr(err)
	}

	response, err := client.Security.PutRole(roleName, &buffer)
	if err != nil {
		return diag.FromErr(err)
	}

	defer response.Body.Close()

	if response.IsError() {
		return diag.Errorf("Failed to create role: [%d] %s", response.StatusCode, response.String())
	}

	// Empty the response body...
	io.Copy(ioutil.Discard, response.Body)

	data.SetId(roleName)

	return resourceRoleRead(context, data, state)
}

func resourceRoleRead(context context.Context, data *schema.ResourceData, state interface{}) diag.Diagnostics {
	client := state.(*api.Client)

	var diags diag.Diagnostics

	name := data.Id()

	response, err := client.Security.GetRole(client.Security.GetRole.WithName(name))

	if err != nil {
		return diag.FromErr(err)
	}

	defer response.Body.Close()

	if response.IsError() {
		return diag.Errorf("Failed to read role: [%d] %s", response.StatusCode, response.String())
	}

	var roleResponse map[string]roleModel

	err = json.NewDecoder(response.Body).Decode(&roleResponse)

	if err != nil {
		return diag.FromErr(err)
	}

	role, exists := roleResponse[name]

	if !exists {
		return diag.Errorf("Role does not exist!")
	}

	if err == nil {
		data.Set(applicationsKey, role.Applications)
	}

	if err == nil {
		data.Set(clusterKey, role.Cluster)
	}

	if err == nil {
		data.Set(indicesKey, role.Indices)
	}

	if err == nil {
		data.Set(metadataKey, role.Metadata)
	}

	if err == nil {
		data.Set(runAsKey, role.RunAs)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRoleDelete(context context.Context, data *schema.ResourceData, state interface{}) diag.Diagnostics {
	client := state.(*api.Client)

	var diags diag.Diagnostics

	name := data.Id()

	response, err := client.Security.DeleteRole(name)

	if err != nil {
		return diag.FromErr(err)
	}

	defer response.Body.Close()

	if response.IsError() {
		return diag.Errorf("Failed to delete role: [%d] %s", response.StatusCode, response.String())
	}

	io.Copy(ioutil.Discard, response.Body)

	return diags
}

func mapApplications(source []interface{}) []applicationModel {
	var result []applicationModel
	for _, item := range source {
		itemMap := item.(map[string]interface{})

		application := applicationModel{
			Application: itemMap[nameKey].(string),
		}

		if privileges, ok := itemMap[privilegesKey]; ok {
			application.Privileges = mapStringArray(privileges.([]interface{}))
		}

		if resources, ok := itemMap[resourcesKey]; ok {
			application.Resources = mapStringArray(resources.([]interface{}))
		}

		result = append(result, application)
	}
	return result
}

func mapIndices(source []interface{}) []indexModel {
	var result []indexModel
	for _, item := range source {
		itemMap := item.(map[string]interface{})

		index := indexModel{
			Names:      mapStringArray(itemMap[namesKey].([]interface{})),
			Privileges: mapStringArray(itemMap[privilegesKey].([]interface{})),
		}

		if query, ok := itemMap[queryKey]; ok {
			index.Query = query.(string)
		}

		if allowUnrestrictedIndices, ok := itemMap[allowUnRestrictedIndicesKey]; ok {
			index.AllowUnRestrictedIndices = allowUnrestrictedIndices.(bool)
		}

		if fieldSecuritySource, ok := itemMap[fieldSecurityKey]; ok {
			sourceList := fieldSecuritySource.([]interface{})
			if len(sourceList) > 0 {
				sourceMap := sourceList[0].(map[string]interface{})
				fieldSecurity := fieldSecurityModel{
					Grant: mapStringArray(sourceMap[grantKey].([]interface{})),
				}
				index.FieldSecurity = fieldSecurity
			}
		}

		result = append(result, index)
	}
	return result
}
