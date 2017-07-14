package commands

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

type RegenerateCommand struct {
	CredentialIdentifier string `required:"yes" short:"n" long:"name" description:"Selects the credential to regenerate"`
	OutputJson           bool   `long:"output-json" description:"Return response in JSON format"`
}

func (cmd RegenerateCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	repository := repositories.NewCredentialRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, &cfg)

	credential, err := action.DoAction(client.NewRegenerateCredentialRequest(cfg, cmd.CredentialIdentifier), cmd.CredentialIdentifier)
	if err != nil {
		return err
	}

	models.Println(credential, cmd.OutputJson)

	return nil
}
