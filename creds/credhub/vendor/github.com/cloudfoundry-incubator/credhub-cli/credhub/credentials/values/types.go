// CredHub credential value types
package values

type Value string

type JSON interface{}

type Password string

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Certificate struct {
	Ca          string `json:"ca"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key" yaml:"private_key"`
}

type RSA struct {
	PublicKey  string `json:"public_key" yaml:"public_key"`
	PrivateKey string `json:"private_key" yaml:"private_key"`
}

type SSH struct {
	PublicKey  string `json:"public_key" yaml:"public_key"`
	PrivateKey string `json:"private_key" yaml:"private_key"`
}
