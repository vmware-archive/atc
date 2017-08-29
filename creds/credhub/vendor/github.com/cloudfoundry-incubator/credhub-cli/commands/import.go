package commands

import (
	"fmt"

	"net/http"

	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

type ImportCommand struct {
	File string `short:"f" long:"file" description:"File containing credentials to import. File must be in yaml format containing a list of credentials under the key 'credentials'. Name, type and value are required for each credential in the list." required:"true"`
}

var (
	err        error
	repository repositories.Repository
	bulkImport models.CredentialBulkImport
	request    *http.Request
)

func (cmd ImportCommand) Execute([]string) error {
	err = bulkImport.ReadFile(cmd.File)

	if err != nil {
		return err
	}

	setCredentials(bulkImport)

	return nil
}

func setCredentials(bulkImport models.CredentialBulkImport) {
	var name string

	cfg := config.ReadConfig()
	repository = repositories.NewCredentialRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, &cfg)

	for _, credential := range bulkImport.Credentials {
		request = client.NewSetRequest(cfg, credential)

		switch credentialName := credential["name"].(type) {
		case string:
			name = credentialName
		default:
			name = ""
		}

		result, err := action.DoAction(request, name)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}

		models.Println(result, false)
	}
}
