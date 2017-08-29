package config

import "github.com/cloudfoundry-incubator/credhub-cli/errors"

func ValidateConfig(c Config) error {
	if c.ApiURL == "" {
		return errors.NewNoTargetUrlError()
	} else if c.AccessToken == "" || c.AccessToken == "revoked" {
		return errors.NewRevokedTokenError()
	}

	return nil
}
