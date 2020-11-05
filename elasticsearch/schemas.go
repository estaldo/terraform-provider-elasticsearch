package elasticsearch

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var metadataSchema = schema.Schema{
	Type:        schema.TypeMap,
	Optional:    true,
	Description: "Arbitrary metadata that you want to associate with the user.",
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
}

var userResource = schema.Resource{
	Schema: map[string]*schema.Schema{
		usernameKey: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "An identifier for the user.",
		},
		passwordKey: {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The user’s password. Passwords must be at least 6 characters long.",
		},
		emailKey: {
			Type:        schema.TypeString,
			Sensitive:   true,
			Optional:    true,
			Description: "The email of the user.",
		},
		enabledKey: {
			Type:        schema.TypeBool,
			Default:     true,
			Optional:    true,
			Description: "Specifies whether the user is enabled.",
		},
		fullNameKey: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The full name of the user.",
		},
		rolesKey: {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A set of roles the user has. The roles determine the user’s access permissions.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		metadataKey: &metadataSchema,
	},
}

var applicationResource = schema.Resource{
	Schema: map[string]*schema.Schema{
		nameKey: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the application to which this entry applies",
		},
		privilegesKey: {
			Type:        schema.TypeList,
			Required:    true,
			Description: "A list of strings, where each element is the name of an application privilege or action.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		resourcesKey: {
			Type:        schema.TypeList,
			Required:    true,
			Description: "A list resources to which the privileges are applied.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	},
}

var fieldSecurityResource = schema.Resource{
	Schema: map[string]*schema.Schema{
		grantKey: {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	},
}

var indexResource = schema.Resource{
	Schema: map[string]*schema.Schema{
		namesKey: {
			Type:        schema.TypeList,
			Required:    true,
			Description: "A list of indices (or index name patterns) to which the permissions in this entry apply.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		privilegesKey: {
			Type:        schema.TypeList,
			Required:    true,
			Description: "The index level privileges that the owners of the role have on the specified indices.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		queryKey: {
			Type:     schema.TypeString,
			Optional: true,
			Description: `
			A search query that defines the documents the owners of the role have read access to. 
			A document within the specified indices must match this query in order for it to be accessible by the owners of the role.`,
		},
		fieldSecurityKey: {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Description: `The document fields that the owners of the role have read access to. 
			For more information, see https://www.elastic.co/guide/en/elasticsearch/reference/current/field-and-document-access-control.html`,
			Elem: &fieldSecurityResource,
		},
		allowUnRestrictedIndicesKey: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	},
}

var roleResource = schema.Resource{
	Schema: map[string]*schema.Schema{
		nameKey: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the role.",
		},
		applicationsKey: {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &applicationResource,
			Description: "A list of application privilege entries.",
		},
		clusterKey: {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: `A list of cluster privileges. 
			These privileges define the cluster level actions that users with this role are able to execute.`,
		},
		indicesKey: {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &indexResource,
			Description: "A list of indices permissions entries.",
		},
		metadataKey: {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Optional meta-data. Within the metadata object, keys that begin with _ are reserved for system usage.",
		},
		runAsKey: {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: `A list of users that the owners of this role can impersonate. 
			For more information, see https://www.elastic.co/guide/en/elasticsearch/reference/current/run-as-privilege.html`,
		},
	},
}

var apiKeyResource = schema.Resource{
	Schema: map[string]*schema.Schema{
		nameKey: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Specifies the name for this API key.",
		},
		expirationKey: {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Expiration time for the API key. By default, API keys never expire.",
		},
		apiKeyKey: {
			Type:     schema.TypeString,
			Computed: true,
		},
		roleDescriptorsKey: {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Elem:     &roleResource,
			Description: `An array of role descriptors for this API key. 
			This parameter is optional. When it is not specified or is an empty array, 
			then the API key will have a point in time snapshot of permissions of the authenticated user. 
			If you supply role descriptors then the resultant permissions would be an intersection of 
			API keys permissions and authenticated user’s permissions thereby limiting the access scope for API keys.
			The structure of role descriptor is the same as the request for create role API. 
			For more details, see https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-put-role.html.`,
		},
	},
}
