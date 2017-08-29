package actions

import (
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type ServerInfo struct {
	httpClient client.HttpClient
	config     config.Config
}
