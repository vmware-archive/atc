package repositories

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/errors"

	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"

	credhub_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
)

type credentialQueryRepository struct {
	httpClient client.HttpClient
}

func NewCredentialQueryRepository(httpClient client.HttpClient) Repository {
	return credentialQueryRepository{httpClient: httpClient}
}

func (r credentialQueryRepository) SendRequest(request *http.Request, identifier string) (models.Printable, error) {
	credentialResponse := models.CredentialResponse{}
	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return credentialResponse, err
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&credentialResponse.ResponseBody)

	if err != nil {
		return credentialResponse, credhub_errors.NewResponseError()
	}
	if len(credentialResponse.ResponseBody["credentials"].([]interface{})) == 0 {
		return credentialResponse, errors.NewNoMatchingCredentialsFoundError()
	}

	return credentialResponse, nil
}
