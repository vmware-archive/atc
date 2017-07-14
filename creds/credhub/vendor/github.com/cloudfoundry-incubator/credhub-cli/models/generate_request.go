package models

type GenerateRequest struct {
	Name           string                `json:"name"`
	CredentialType string                `json:"type"`
	Overwrite      *bool                 `json:"overwrite"`
	Parameters     *GenerationParameters `json:"parameters"`
	Value          *ProvidedValue        `json:"value,omitempty"`
}

type GenerationParameters struct {
	IncludeSpecial   bool     `json:"include_special,omitempty"`
	ExcludeNumber    bool     `json:"exclude_number,omitempty"`
	ExcludeUpper     bool     `json:"exclude_upper,omitempty"`
	ExcludeLower     bool     `json:"exclude_lower,omitempty"`
	Length           int      `json:"length,omitempty"`
	CommonName       string   `json:"common_name,omitempty"`
	Organization     string   `json:"organization,omitempty"`
	OrganizationUnit string   `json:"organization_unit,omitempty"`
	Locality         string   `json:"locality,omitempty"`
	State            string   `json:"state,omitempty"`
	Country          string   `json:"country,omitempty"`
	AlternativeName  []string `json:"alternative_names,omitempty"`
	ExtendedKeyUsage []string `json:"extended_key_usage,omitempty"`
	KeyUsage         []string `json:"key_usage,omitempty"`
	KeyLength        int      `json:"key_length,omitempty"`
	Duration         int      `json:"duration,omitempty"`
	Ca               string   `json:"ca,omitempty"`
	SelfSign         bool     `json:"self_sign,omitempty"`
	IsCA             bool     `json:"is_ca,omitempty"`
	SshComment       string   `json:"ssh_comment,omitempty"`
}

type ProvidedValue struct {
	Username string `json:"username,omitempty"`
}
