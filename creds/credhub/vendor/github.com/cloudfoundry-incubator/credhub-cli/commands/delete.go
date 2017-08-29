package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

type DeleteCommand struct {
	CredentialIdentifier string `short:"n" long:"name" required:"yes" description:"Name of the credential to delete"`
}

func (cmd DeleteCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	repository := repositories.NewCredentialRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, &cfg)

	_, err := action.DoAction(client.NewDeleteCredentialRequest(cfg, cmd.CredentialIdentifier), cmd.CredentialIdentifier)

	if err == nil {
		fmt.Println("Credential successfully deleted")
	}

	return err
}
