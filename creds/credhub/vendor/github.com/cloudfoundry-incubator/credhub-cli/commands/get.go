package commands

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

type GetCommand struct {
	Name       string `short:"n" long:"name" description:"Name of the credential to retrieve"`
	Id         string `long:"id" description:"ID of the credential to retrieve"`
	OutputJson bool   `long:"output-json" description:"Return response in JSON format"`
}

func (cmd GetCommand) Execute([]string) error {
	var (
		credential models.Printable
		err        error
	)

	cfg := config.ReadConfig()
	repository := repositories.NewCredentialRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, &cfg)

	if cmd.Name != "" {
		credential, err = action.DoAction(client.NewGetCredentialByNameRequest(cfg, cmd.Name), cmd.Name)
	} else if cmd.Id != "" {
		credential, err = action.DoAction(client.NewGetCredentialByIdRequest(cfg, cmd.Id), cmd.Id)
	} else {
		return errors.NewMissingGetParametersError()
	}

	if err != nil {
		return err
	}

	models.Println(credential, cmd.OutputJson)

	return nil
}
