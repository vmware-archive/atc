package commands

import (
	"fmt"

	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/fatih/color"
)

var warning = color.New(color.Bold, color.FgYellow).PrintlnFunc()
var deprecation = color.New(color.Bold, color.FgRed).PrintlnFunc()

type ApiCommand struct {
	Server            ApiPositionalArgs `positional-args:"yes"`
	ServerFlagUrl     string            `short:"s" long:"server" description:"URI of API server to target"`
	CaCert            []string          `long:"ca-cert" description:"Trusted CA for API and UAA TLS connections"`
	SkipTlsValidation bool              `long:"skip-tls-validation" description:"Skip certificate validation of the API endpoint. Not recommended!"`
}

type ApiPositionalArgs struct {
	ServerUrl string `positional-arg-name:"SERVER" description:"URI of API server to target"`
}

func (cmd ApiCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	serverUrl := targetUrl(cmd)

	cfg.CaCert = cmd.CaCert

	if serverUrl == "" {
		if cfg.ApiURL != "" {
			fmt.Println(cfg.ApiURL)
		} else {
			return errors.NewNoTargetUrlError()
		}
	} else {
		existingCfg := cfg
		err := GetApiInfo(&cfg, serverUrl, cmd.SkipTlsValidation)
		if err != nil {
			return err
		}

		fmt.Println("Setting the target url:", cfg.ApiURL)

		if existingCfg.AuthURL != cfg.AuthURL {
			RevokeTokenIfNecessary(existingCfg)
			MarkTokensAsRevokedInConfig(&cfg)
		}
		config.WriteConfig(cfg)
	}

	return nil
}

func GetApiInfo(cfg *config.Config, serverUrl string, skipTlsValidation bool) error {
	serverUrl = AddDefaultSchemeIfNecessary(serverUrl)
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		return err
	}

	cfg.ApiURL = parsedUrl.String()

	cfg.InsecureSkipVerify = skipTlsValidation
	credhubInfo, err := actions.NewInfo(client.NewHttpClient(*cfg), *cfg).GetServerInfo()
	if err != nil {
		return err
	}
	cfg.AuthURL = credhubInfo.AuthServer.Url

	if parsedUrl.Scheme != "https" {
		warning("Warning: Insecure HTTP API detected. Data sent to this API could be intercepted" +
			" in transit by third parties. Secure HTTPS API endpoints are recommended.")
	} else {
		if skipTlsValidation {
			warning("Warning: The targeted TLS certificate has not been verified for this connection.")
			deprecation("Warning: The --skip-tls-validation flag is deprecated. Please use --ca-cert instead.")
		}
	}

	return nil
}

func targetUrl(cmd ApiCommand) string {
	if cmd.Server.ServerUrl != "" {
		return cmd.Server.ServerUrl
	} else {
		return cmd.ServerFlagUrl
	}
}
