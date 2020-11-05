package elasticsearch

type userModel struct {
	Username string            `json:"username"`
	Password string            `json:"password"`
	Email    string            `json:"email"`
	Enabled  bool              `json:"enabled"`
	FullName string            `json:"full_name"`
	Roles    []string          `json:"roles"`
	Metadata map[string]string `json:"metadata"`
}

type fieldSecurityModel struct {
	Grant []string `json:"grant"`
}

type indexModel struct {
	Names                    []string           `json:"names"`
	Privileges               []string           `json:"privileges"`
	FieldSecurity            fieldSecurityModel `json:"field_security,omitempty"`
	Query                    string             `json:"query,omitempty"`
	AllowUnRestrictedIndices bool               `json:"allow_restricted_indices,omitempty"`
}

type applicationModel struct {
	Application string   `json:"application"`
	Privileges  []string `json:"privileges"`
	Resources   []string `json:"resources"`
}

type roleModel struct {
	Cluster      []string           `json:"cluster"`
	Indices      []indexModel       `json:"indices,omitempty"`
	Applications []applicationModel `json:"applications,omitempty"`
	RunAs        []string           `json:"run_as,omitempty"`
	Metadata     map[string]string  `json:"metadata,omitempty"`
}

type apiKeyModel struct {
	Name           string               `json:"name"`
	Expiration     string               `json:"expiration,omitempty"`
	RoleDescriptor map[string]roleModel `json:"role_descriptors,omitempty"`
}

type apiKeyCreateResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Expiration int64  `json:"expiration"`
	APIKey     string `json:"api_key"`
}

type apiKeyGetResponse struct {
	APIKeys []apiKeyCreateResponse `json:"api_keys"`
}
