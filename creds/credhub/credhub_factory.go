package credhub

import (
	"sync"
	"time"

	"code.cloudfoundry.org/lager"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/commands"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
	"github.com/concourse/atc/creds"
)

type credhubFactory struct {
	config *config.Config

	prefix string

	configL *sync.RWMutex

	loggedIn chan struct{}
}

func NewCredhubFactory(logger lager.Logger, config config.Config, auth AuthConfig, prefix string) *credhubFactory {
	factory := &credhubFactory{
		config: &config,

		prefix: prefix,

		configL:  new(sync.RWMutex),
		loggedIn: make(chan struct{}),
	}

	go factory.authLoop(logger, auth)

	return factory
}

func (factory *credhubFactory) NewVariables(teamName string, pipelineName string) creds.Variables {
	<-factory.loggedIn

	return &Credhub{
		Config: factory.currentConfig(),

		PathPrefix:   factory.prefix,
		TeamName:     teamName,
		PipelineName: pipelineName,
	}
}

func (factory *credhubFactory) currentConfig() *config.Config {
	factory.configL.RLock()
	config := factory.config
	factory.configL.RUnlock()
	return config
}

func (factory *credhubFactory) authLoop(logger lager.Logger, auth AuthConfig) {
	for {
		currentConfig := factory.currentConfig()

		var cfg *config.Config
		var delay time.Duration
		if currentConfig.AccessToken == "" {
			cfg, delay = factory.login(logger.Session("login"), auth)
		} else {
			cfg, delay = factory.renew(logger.Session("renew"))
		}

		if cfg.AccessToken != "" {
			factory.configL.Lock()
			factory.config = cfg
			if currentConfig.AccessToken == "" {
				close(factory.loggedIn)
			}
			factory.configL.Unlock()
		}

		time.Sleep(delay)
	}
}

func (factory *credhubFactory) login(logger lager.Logger, auth AuthConfig) (*config.Config, time.Duration) {
	configCopy := *factory.config
	if configCopy.AccessToken != "" {
		return factory.config, 0
	}

	err := commands.GetApiInfo(&configCopy, configCopy.ApiURL, configCopy.InsecureSkipVerify)
	if err != nil {
		logger.Error("failed getting Credhub api info", err)
		return factory.config, time.Second
	}

	httpClient := client.NewHttpClient(configCopy)
	token, err := actions.NewAuthToken(httpClient, configCopy).GetAuthTokenByClientCredential(
		auth.ClientName,
		auth.ClientSecret,
	)
	if err != nil {
		logger.Error("failed getting auth token from UAA", err)
		return factory.config, time.Second
	}

	logger.Info("succeeded", lager.Data{
		"token-type":     token.TokenType,
		"lease-duration": token.ExpiresIn,
	})

	configCopy.AccessToken = token.AccessToken
	configCopy.RefreshToken = token.RefreshToken

	return &configCopy, (time.Duration(token.ExpiresIn) * time.Second) / 2
}

func (factory *credhubFactory) renew(logger lager.Logger) (*config.Config, time.Duration) {
	configCopy := *factory.config

	if configCopy.RefreshToken == "" {
		configCopy.AccessToken = ""
		return &configCopy, time.Second
	}

	refreshRequest := client.NewRefreshTokenRequest(configCopy)
	repository := repositories.NewAuthRepository(client.NewHttpClient(configCopy), true)
	refreshedToken, err := repository.SendRequest(refreshRequest, "")

	if err != nil {
		logger.Error("failed refreshing UAA auth token", err)
		return &configCopy, time.Second
	}

	logger.Info("succeeded", lager.Data{
		"token-type":     refreshedToken.(models.Token).TokenType,
		"lease-duration": refreshedToken.(models.Token).ExpiresIn,
	})

	configCopy.AccessToken = refreshedToken.(models.Token).AccessToken
	configCopy.RefreshToken = refreshedToken.(models.Token).RefreshToken

	return &configCopy, (time.Duration(refreshedToken.(models.Token).ExpiresIn) * time.Second) / 2
}
