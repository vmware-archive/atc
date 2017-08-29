// CredHub credential types
package credentials

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

// Base fields of a credential
type Base struct {
	Name             string `json:"name"`
	VersionCreatedAt string `json:"version_created_at" yaml:"version_created_at"`
}

type Metadata struct {
	Base `yaml:",inline"`
	Id   string `json:"id"`
	Type string `json:"type"`
}

// A generic credential
//
// Used when the Type of the credential is not known ahead of time.
//
// Value will be as unmarshalled by https://golang.org/pkg/encoding/json/#Unmarshal
type Credential struct {
	Metadata `yaml:",inline"`
	Value    interface{} `json:"value"`
}

// A Value type credential
type Value struct {
	Metadata `yaml:",inline"`
	Value    values.Value `json:"value"`
}

// A JSON type credential
type JSON struct {
	Metadata
	Value json.RawMessage `json:"value"`
}

func (j JSON) MarshalYAML() (interface{}, error) {
	var x interface{}

	json.Unmarshal(j.Value, &x)

	return struct {
		Metadata `yaml:",inline"`
		Value    interface{}
	}{
		Metadata: j.Metadata,
		Value:    x,
	}, nil

}

// A Password type credential
type Password struct {
	Metadata `yaml:",inline"`
	Value    values.Password `json:"value"`
}

// A User type credential
type User struct {
	Metadata `yaml:",inline"`
	Value    struct {
		values.User  `yaml:",inline"`
		PasswordHash string `json:"password_hash" yaml:"password_hash"`
	} `json:"value"`
}

// A Certificate type credential
type Certificate struct {
	Metadata `yaml:",inline"`
	Value    values.Certificate `json:"value"`
}

// An RSA type credential
type RSA struct {
	Metadata `yaml:",inline"`
	Value    values.RSA `json:"value"`
}

// An SSH type credential
type SSH struct {
	Metadata `yaml:",inline"`
	Value    struct {
		values.SSH           `yaml:",inline"`
		PublicKeyFingerprint string `json:"public_key_fingerprint" yaml:"public_key_fingerprint"`
	} `json:"value"`
}
