package credhub

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

// GetById returns a credential version by ID. The returned credential may be of any type.
func (ch *CredHub) GetById(id string) (credentials.Credential, error) {
	panic("Not implemented")
}

// GetAll returns all credential versions for a given credential name. The returned credential may be of any type.
func (ch *CredHub) GetAll(name string) ([]credentials.Credential, error) {
	panic("Not implemented")
}

// Get returns the current credential version for a given credential name. The returned credential may be of any type.
func (ch *CredHub) Get(name string) (credentials.Credential, error) {
	var cred credentials.Credential
	err := ch.getCurrentCredential(name, &cred)
	return cred, err
}

// GetValue returns the current credential version for a given credential name. The returned credential must be of type 'value'.
func (ch *CredHub) GetValue(name string) (credentials.Value, error) {
	var cred credentials.Value
	err := ch.getCurrentCredential(name, &cred)

	return cred, err
}

// GetJSON returns the current credential version for a given credential name. The returned credential must be of type 'json'.
func (ch *CredHub) GetJSON(name string) (credentials.JSON, error) {
	var cred credentials.JSON
	err := ch.getCurrentCredential(name, &cred)

	return cred, err
}

// GetPassword returns the current credential version for a given credential name. The returned credential must be of type 'password'.
func (ch *CredHub) GetPassword(name string) (credentials.Password, error) {
	var cred credentials.Password
	err := ch.getCurrentCredential(name, &cred)

	return cred, err
}

// GetUser returns the current credential version for a given credential name. The returned credential must be of type 'user'.
func (ch *CredHub) GetUser(name string) (credentials.User, error) {
	var cred credentials.User
	err := ch.getCurrentCredential(name, &cred)

	return cred, err
}

// GetCertificate returns the current credential version for a given credential name. The returned credential must be of type 'certificate'.
func (ch *CredHub) GetCertificate(name string) (credentials.Certificate, error) {
	var cred credentials.Certificate
	err := ch.getCurrentCredential(name, &cred)

	return cred, err
}

// GetRSA returns the current credential version for a given credential name. The returned credential must be of type 'rsa'.
func (ch *CredHub) GetRSA(name string) (credentials.RSA, error) {
	var cred credentials.RSA
	err := ch.getCurrentCredential(name, &cred)

	return cred, err
}

// GetSSH returns the current credential version for a given credential name. The returned credential must be of type 'ssh'.
func (ch *CredHub) GetSSH(name string) (credentials.SSH, error) {
	var cred credentials.SSH
	err := ch.getCurrentCredential(name, &cred)

	return cred, err
}

func (ch *CredHub) getCurrentCredential(name string, cred interface{}) error {
	query := url.Values{}
	query.Set("name", name)
	query.Set("current", "true")

	resp, err := ch.Request(http.MethodGet, "/api/v1/data", query, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	response := make(map[string][]json.RawMessage)

	if err := dec.Decode(&response); err != nil {
		return err
	}

	var ok bool
	var data []json.RawMessage

	if data, ok = response["data"]; !ok || len(data) == 0 {
		return errors.New("response did not contain any credentials")
	}

	rawMessage := data[0]

	return json.Unmarshal(rawMessage, cred)
}
