package genericoauth

import (
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/auth/verifier"
	"github.com/concourse/atc/db"
	"golang.org/x/oauth2"
)

const ProviderName = "oauth"

type NoopVerifier struct{}

func (v NoopVerifier) Verify(logger lager.Logger, client *http.Client) (bool, error) {
	return true, nil
}

func NewProvider(
	genericOAuth *db.GenericOAuth,
	redirectURL string,
) Provider {
	endpoint := oauth2.Endpoint{}
	if genericOAuth.AuthURL != "" && genericOAuth.TokenURL != "" {
		endpoint.AuthURL = genericOAuth.AuthURL
		endpoint.TokenURL = genericOAuth.TokenURL
	}

	return Provider{
		Verifier: NoopVerifier{},
		Config: &oauth2.Config{
			ClientID:     genericOAuth.ClientID,
			ClientSecret: genericOAuth.ClientSecret,
			Endpoint:     endpoint,
			RedirectURL:  redirectURL,
		},
		ConfiguredDisplayName: genericOAuth.DisplayName,
	}
}

type Provider struct {
	*oauth2.Config
	// oauth2.Config implements the required Provider methods:
	// AuthCodeURL(string, ...oauth2.AuthCodeOption) string
	// Exchange(context.Context, string) (*oauth2.Token, error)
	// Client(context.Context, *oauth2.Token) *http.Client

	verifier.Verifier
	ConfiguredDisplayName string
}

func (provider Provider) DisplayName() string {
	return provider.ConfiguredDisplayName
}

func (Provider) PreTokenClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
}
