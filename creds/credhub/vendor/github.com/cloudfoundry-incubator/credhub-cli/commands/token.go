package commands

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func init() {
	CredHub.Token = func() {
		cfg := config.ReadConfig()
		if cfg.AccessToken != "" && cfg.AccessToken != "revoked" {

			refresh_request := client.NewRefreshTokenRequest(cfg)
			repository := repositories.NewAuthRepository(client.NewHttpClient(cfg), true)
			refreshed_token, err := repository.SendRequest(refresh_request, "")

			if err != nil {
				fmt.Println("Bearer " + cfg.AccessToken)
			}

			cfg.AccessToken = refreshed_token.(models.Token).AccessToken
			cfg.RefreshToken = refreshed_token.(models.Token).RefreshToken

			config.WriteConfig(cfg)

			fmt.Println("Bearer " + cfg.AccessToken)
		} else {
			fmt.Fprint(os.Stderr, "You are not currently authenticated. Please log in to continue.")
		}
		os.Exit(0)
	}
}
