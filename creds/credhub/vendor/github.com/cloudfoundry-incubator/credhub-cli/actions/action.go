package actions

import (
	"net/http"

	"reflect"

	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

type Action struct {
	repository     repositories.Repository
	config         config.Config
	AuthRepository repositories.Repository
}

func NewAction(repository repositories.Repository, cfg *config.Config) Action {
	err := config.ValidateConfig(*cfg)
	var token models.Token

	if reflect.DeepEqual(err, errors.NewRevokedTokenError()) && (os.Getenv("CREDHUB_CLIENT") != "" || os.Getenv("CREDHUB_SECRET") != "") {
		token, err = NewAuthToken(client.NewHttpClient(*cfg), *cfg).GetAuthTokenByClientCredential(os.Getenv("CREDHUB_CLIENT"), os.Getenv("CREDHUB_SECRET"))
		if err == nil {
			cfg.AccessToken = token.AccessToken
			cfg.RefreshToken = token.RefreshToken
			config.WriteConfig(*cfg)
		}
	}

	action := Action{repository: repository, config: *cfg}
	action.AuthRepository = repositories.NewAuthRepository(client.NewHttpClient(*cfg), true)
	return action
}

func (action Action) DoAction(req *http.Request, identifier string) (models.Printable, error) {
	err := config.ValidateConfig(action.config)

	if err != nil {
		return nil, err
	}

	bodyClone := client.NewBodyClone(req)

	item, err := action.repository.SendRequest(req, identifier)

	if reflect.DeepEqual(err, errors.NewAccessTokenExpiredError()) {
		req.Body = bodyClone
		item, err = action.refreshTokenAndResendRequest(req, identifier)
	}
	return item, err
}

func (action Action) refreshTokenAndResendRequest(req *http.Request, identifier string) (models.Printable, error) {
	err := action.refreshToken()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+action.config.AccessToken)
	item, err := action.repository.SendRequest(req, identifier)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (action *Action) refreshToken() error {
	var (
		token models.Token
		err   error
	)
	if os.Getenv("CREDHUB_CLIENT") != "" || os.Getenv("CREDHUB_SECRET") != "" {
		token, err = NewAuthToken(client.NewHttpClient(action.config), action.config).GetAuthTokenByClientCredential(os.Getenv("CREDHUB_CLIENT"), os.Getenv("CREDHUB_SECRET"))
		if err != nil {
			return err
		}

		action.config.AccessToken = token.AccessToken
		action.config.RefreshToken = token.RefreshToken
	} else {
		refresh_request := client.NewRefreshTokenRequest(action.config)
		refreshed_token, err := action.AuthRepository.SendRequest(refresh_request, "")

		if err != nil {
			return errors.NewRefreshError()
		}

		action.config.AccessToken = refreshed_token.(models.Token).AccessToken
		action.config.RefreshToken = refreshed_token.(models.Token).RefreshToken
	}

	config.WriteConfig(action.config)

	return nil
}
