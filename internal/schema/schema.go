package schema

type CapabilitySchema struct {
	Name            string            `json:"name"`
	ContractVersion string            `json:"contract_version"`
	Description     string            `json:"description"`
	Parameters      []CapabilityParam `json:"parameters,omitempty"`
	Rules           []CapabilityRule  `json:"rules,omitempty"`
	OutputFields    map[string]any    `json:"output_fields,omitempty"`
}

type CapabilityParam struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Required    bool     `json:"required,omitempty"`
	Description string   `json:"description,omitempty"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	ItemsType   string   `json:"items_type,omitempty"`
	Encoding    string   `json:"encoding,omitempty"`
	FlagName    string   `json:"flag_name,omitempty"`
}

type CapabilityRule struct {
	Kind    string   `json:"kind"`
	Params  []string `json:"params,omitempty"`
	Message string   `json:"message,omitempty"`
}
