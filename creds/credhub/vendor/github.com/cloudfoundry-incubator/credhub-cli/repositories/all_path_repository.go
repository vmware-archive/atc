package repositories

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/errors"

	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"

	credhub_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
)

type allPathRepository struct {
	httpClient client.HttpClient
}

func NewAllPathRepository(httpClient client.HttpClient) Repository {
	return allPathRepository{httpClient: httpClient}
}

func (r allPathRepository) SendRequest(request *http.Request, ignoredIdentifier string) (models.Printable, error) {
	credential_paths := models.CredentialResponse{}

	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return credential_paths, err
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&credential_paths.ResponseBody)

	if err != nil {
		return credential_paths, credhub_errors.NewResponseError()
	}

	if len(credential_paths.ResponseBody["paths"].([]interface{})) == 0 {
		return credential_paths, errors.NewNoMatchingCredentialsFoundError()
	}
	return credential_paths, nil
}
