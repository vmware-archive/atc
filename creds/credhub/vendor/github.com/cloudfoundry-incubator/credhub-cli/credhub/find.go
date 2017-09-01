package credhub

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

// FindByPartialName retrieves a list of stored credential names which contain the search.
func (ch *CredHub) FindByPartialName(nameLike string) ([]credentials.Base, error) {
	panic("Not implemented")
}

// FindByPath retrieves a list of stored credential names which are within the specified path.
func (ch *CredHub) FindByPath(path string) ([]credentials.Base, error) {
	var creds map[string][]credentials.Base

	query := url.Values{}
	query.Set("path", path)

	resp, err := ch.Request(http.MethodGet, "/api/v1/data", query, nil)

	if err != nil {
		return []credentials.Base{}, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &creds)

	if err != nil {
		return []credentials.Base{}, err
	}

	return creds["credentials"], nil
}

// ShowAllPaths retrieves a list of all paths which contain credentials.
func (ch *CredHub) ShowAllPaths() ([]credentials.Path, error) {
	panic("Not implemented")
}
