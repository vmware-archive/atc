package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	credhub_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type credentialRepository struct {
	httpClient client.HttpClient
}

func NewCredentialRepository(httpClient client.HttpClient) Repository {
	return credentialRepository{httpClient: httpClient}
}

func (r credentialRepository) SendRequest(request *http.Request, identifier string) (models.Printable, error) {
	credentialResponse := models.CredentialResponse{}
	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return credentialResponse, err
	}

	if request.Method == "DELETE" {
		return credentialResponse, nil
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&credentialResponse.ResponseBody)

	if err != nil {
		return credentialResponse, credhub_errors.NewResponseError()
	}

	if data, ok := credentialResponse.ResponseBody["data"].([]interface{}); ok {
		credentialResponse.ResponseBody = data[0].(map[string]interface{})
	}

	return credentialResponse, nil
}
