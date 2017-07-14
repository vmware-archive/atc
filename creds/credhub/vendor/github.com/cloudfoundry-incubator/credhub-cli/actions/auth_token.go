package actions

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

func NewAuthToken(httpClient client.HttpClient, config config.Config) ServerInfo {
	return ServerInfo{httpClient: httpClient, config: config}
}

func (serverInfo ServerInfo) GetAuthTokenByPasswordGrant(user string, pass string) (models.Token, error) {
	request := client.NewPasswordGrantTokenRequest(serverInfo.config, user, pass)
	return getAuthToken(serverInfo, request)
}

func (serverInfo ServerInfo) GetAuthTokenByClientCredential(clientName string, clientSecret string) (models.Token, error) {
	request := client.NewClientCredentialsGrantTokenRequest(serverInfo.config, clientName, clientSecret)
	return getAuthToken(serverInfo, request)
}

func getAuthToken(serverInfo ServerInfo, request *http.Request) (models.Token, error) {
	response, err := serverInfo.httpClient.Do(request)
	if err != nil {
		return models.Token{}, errors.NewNetworkError(err)
	}

	if response.StatusCode != http.StatusOK {
		return models.Token{}, errors.NewAuthorizationError()
	}

	token := new(models.Token)

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(token)

	if err != nil {
		return models.Token{}, err
	}

	return *token, nil
}
