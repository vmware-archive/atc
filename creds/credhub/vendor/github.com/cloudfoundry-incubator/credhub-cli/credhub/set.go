package credhub

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

// SetValue sets a value credential with a user-provided value.
func (ch *CredHub) SetValue(name string, value values.Value, overwrite bool) (credentials.Value, error) {
	var cred credentials.Value
	err := ch.setCredential(name, "value", value, overwrite, &cred)

	return cred, err
}

// SetJSON sets a JSON credential with a user-provided value.
func (ch *CredHub) SetJSON(name string, value values.JSON, overwrite bool) (credentials.JSON, error) {
	var cred credentials.JSON
	err := ch.setCredential(name, "json", value, overwrite, &cred)

	return cred, err
}

// SetPassword sets a password credential with a user-provided value.
func (ch *CredHub) SetPassword(name string, value values.Password, overwrite bool) (credentials.Password, error) {
	var cred credentials.Password
	err := ch.setCredential(name, "password", value, overwrite, &cred)

	return cred, err
}

// SetUser sets a user credential with a user-provided value.
func (ch *CredHub) SetUser(name string, value values.User, overwrite bool) (credentials.User, error) {
	var cred credentials.User
	err := ch.setCredential(name, "user", value, overwrite, &cred)

	return cred, err
}

// SetCertificate sets a certificate credential with a user-provided value.
func (ch *CredHub) SetCertificate(name string, value values.Certificate, overwrite bool) (credentials.Certificate, error) {
	var cred credentials.Certificate
	err := ch.setCredential(name, "certificate", value, overwrite, &cred)

	return cred, err
}

// SetRSA sets an RSA credential with a user-provided value.
func (ch *CredHub) SetRSA(name string, value values.RSA, overwrite bool) (credentials.RSA, error) {
	var cred credentials.RSA
	err := ch.setCredential(name, "rsa", value, overwrite, &cred)

	return cred, err
}

// SetSSH sets an SSH credential with a user-provided value.
func (ch *CredHub) SetSSH(name string, value values.SSH, overwrite bool) (credentials.SSH, error) {
	var cred credentials.SSH
	err := ch.setCredential(name, "ssh", value, overwrite, &cred)

	return cred, err
}

func (ch *CredHub) setCredential(name, credType string, value interface{}, overwrite bool, cred interface{}) error {
	requestBody := map[string]interface{}{}
	requestBody["name"] = name
	requestBody["type"] = credType
	requestBody["value"] = value
	requestBody["overwrite"] = overwrite
	resp, err := ch.Request(http.MethodPut, "/api/v1/data", nil, requestBody)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(cred)
}
