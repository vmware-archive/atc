package credhub

import (
	"errors"
	"fmt"
	"net/url"

	"code.cloudfoundry.org/lager"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/concourse/atc/creds"
)

type CredhubManager struct {
	URL string `long:"url" description:"Credhub server address used to access secrets."`

	PathPrefix string `long:"path-prefix" default:"/concourse" description:"Path under which to namespace credential lookup."`

	TLS struct {
		CACert   string `long:"ca-cert"              description:"Path to a PEM-encoded CA cert file to use to verify the credhub / UAA server SSL certs."`
		Insecure bool   `long:"insecure-skip-verify" description:"Enable insecure SSL verification."`
	}

	Auth AuthConfig
}

type AuthConfig struct {
	ClientName   string `long:"client-name" description:"Client name for UAA client grant"`
	ClientSecret string `long:"client-secret" description:"Client secret for UAA client grant"`
}

func (manager CredhubManager) IsConfigured() bool {
	return manager.URL != ""
}

func (manager CredhubManager) Validate() error {
	_, err := url.Parse(manager.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %s", err)
	}

	if manager.Auth.ClientName == "" {
		return errors.New("Missing UAA client name for credhub auth")
	}

	if manager.Auth.ClientSecret == "" {
		return errors.New("Missing UAA client secret for credhub auth")
	}

	return nil
}

func (manager CredhubManager) NewVariablesFactory(logger lager.Logger) (creds.VariablesFactory, error) {
	cfg := config.Config{
		ApiURL:             manager.URL,
		InsecureSkipVerify: manager.TLS.Insecure,
		CaCert:             []string{manager.TLS.CACert},
	}

	return NewCredhubFactory(logger, cfg, manager.Auth, manager.PathPrefix), nil
}
