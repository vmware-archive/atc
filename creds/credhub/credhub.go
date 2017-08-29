package credhub

import (
	"path"
	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
	"github.com/cloudfoundry/bosh-cli/director/template"
)

type Credhub struct {
	Config *config.Config

	PathPrefix   string
	TeamName     string
	PipelineName string
}

type Credential struct {
	Value map[string]interface{} `json:"value"`
}

func (c Credhub) Get(varDef template.VariableDefinition) (interface{}, bool, error) {
	var credential Credential
	var found bool
	var err error

	if c.PipelineName != "" {
		credential, found, err = c.findCredential(c.path(c.TeamName, c.PipelineName, varDef.Name))
		if err != nil {
			return nil, false, err
		}
	}

	if !found {
		credential, found, err = c.findCredential(c.path(c.TeamName, varDef.Name))
		if err != nil {
			return nil, false, err
		}
	}

	if !found {
		return nil, false, nil
	}

	return credential.Value, true, nil
}

func (c Credhub) findCredential(path string) (Credential, bool, error) {
	rawCredential, err := actions.NewAction(
		repositories.NewCredentialRepository(client.NewHttpClient(*c.Config)),
		c.Config,
	).DoAction(client.NewGetCredentialByNameRequest(*c.Config, path), path)

	if err != nil {
		return Credential{}, false, err
	}

	if rawCredential != nil {
		var credential Credential
		if err := json.Unmarshal([]byte(rawCredential.ToJson()), &credential); err != nil {
			return Credential{}, false, err
		}
		return credential, true, nil
	}

	return Credential{}, false, nil
}

func (c Credhub) path(segments ...string) string {
	return path.Join(append([]string{c.PathPrefix}, segments...)...)
}

func (c Credhub) List() ([]template.VariableDefinition, error) {
	return []template.VariableDefinition{}, nil
}
