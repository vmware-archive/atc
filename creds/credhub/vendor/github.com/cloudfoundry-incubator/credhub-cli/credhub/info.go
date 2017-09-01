package credhub

import (
	"encoding/json"
	"errors"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
)

// Info returns the targeted CredHub server information.
func (ch *CredHub) Info() (*server.Info, error) {
	response, err := ch.request(ch.Client(), "GET", "/info", nil, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	info := &server.Info{}
	decoder := json.NewDecoder(response.Body)

	if err = decoder.Decode(&info); err != nil {
		return nil, err
	}

	return info, nil
}

// AuthURL returns the targeted CredHub server's trusted authentication server URL.
func (ch *CredHub) AuthURL() (string, error) {
	if ch.authURL != nil {
		return ch.authURL.String(), nil
	}

	info, err := ch.Info()

	if err != nil {
		return "", err
	}

	authUrl := info.AuthServer.URL

	if authUrl == "" {
		return "", errors.New("AuthURL not found")
	}

	return authUrl, nil
}
