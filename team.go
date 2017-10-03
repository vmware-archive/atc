package atc

import "encoding/json"

type Team struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`

	BasicAuth     *BasicAuth     `json:"basic_auth,omitempty"`
	LdapBasicAuth *LdapBasicAuth `json:"ldap_basic_auth,omitempty"`

	Auth map[string]*json.RawMessage `json:"auth,omitempty"`
}

type BasicAuth struct {
	BasicAuthUsername string `json:"basic_auth_username,omitempty"`
	BasicAuthPassword string `json:"basic_auth_password,omitempty"`
}

type LdapBasicAuth struct {
	Server                string `json:"server,omitempty"`
	Port                  uint16 `json:"port,omitempty"`
	TLSEnabled            bool   `json:"tls_enabled,omitempty"`
	TLSInsecureSkipVerify bool   `json:"tls_insecure_skip_verify,omitempty"`
	TLSCA                 string `json:"tls_certificate_authority,omitempty"`
	BindUsername          string `json:"bind_username,omitempty"`
	BindPassword          string `json:"bind_password,omitempty"`
	UserBaseDN            string `json:"user_base_dn,omitempty"`
	GroupBaseDN           string `json:"group_base_dn,omitempty"`
	GroupDN               string `json:"group_dn,omitempty"`
	UserAttribute         string `json:"user_attribute,omitempty"`
	GroupFilter           string `json:"group_filter,omitempty"`
}
