package repositories

import (
	"net/http"

	"encoding/json"
	"errors"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	credhub_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type Repository interface {
	SendRequest(request *http.Request, identifier string) (models.Printable, error)
}

func DoSendRequest(httpClient client.HttpClient, request *http.Request) (*http.Response, error) {
	response, err := httpClient.Do(request)

	if err != nil {
		return nil, credhub_errors.NewNetworkError(err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		decoder := json.NewDecoder(response.Body)
		serverError := models.ServerError{}
		err = decoder.Decode(&serverError)
		if err != nil {
			if response.StatusCode == http.StatusInternalServerError {
				return nil, credhub_errors.NewCatchAllError()
			}
			return nil, err
		}

		if serverError.Error == "access_token_expired" {
			return nil, credhub_errors.NewAccessTokenExpiredError()
		} else if response.StatusCode == http.StatusUnauthorized {
			return nil, errors.New(serverError.ErrorDescription)
		} else if response.StatusCode == http.StatusForbidden {
			return nil, credhub_errors.NewForbiddenError()
		}

		return nil, errors.New(serverError.Error)
	}
	return response, nil
}
