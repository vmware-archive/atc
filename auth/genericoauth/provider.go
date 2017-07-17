package genericoauth

import (
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/concourse/atc"
	"github.com/concourse/atc/auth/provider"
	"github.com/concourse/atc/auth/routes"
	"github.com/concourse/atc/auth/verifier"
	"github.com/hashicorp/go-multierror"
	"github.com/jessevdk/go-flags"
	"github.com/tedsuo/rata"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"net/http"
)

const ProviderName = "oauth"

type Provider struct {
	verifier.Verifier
	Config ConfigOverride
}

type ConfigOverride struct {
	oauth2.Config
	AuthURLParams map[string]string
}

type NoopVerifier struct{}

func init() {
	provider.Register(ProviderName, GenericTeamProvider{})
}

type GenericOAuthConfig struct {
	DisplayName  string `json:"display_name"      long:"display-name"    description:"Name for this auth method on the web UI."`
	ClientID     string `json:"client_id"         long:"client-id"       description:"Application client ID for enabling generic OAuth."`
	ClientSecret string `json:"client_secret"     long:"client-secret"   description:"Application client secret for enabling generic OAuth."`

	AuthURL       string            `json:"auth_url,omitempty"          long:"auth-url"        description:"Generic OAuth provider AuthURL endpoint."`
	AuthURLParams map[string]string `json:"auth_url_params,omitempty"   long:"auth-url-param"  description:"Parameter to pass to the authentication server AuthURL. Can be specified multiple times."`
	Scope         string            `json:"scope,omitempty"             long:"scope"           description:"Optional scope required to authorize user"`
	Authority     string            `json:"authorities,omitempty"       long:"authorities"     description:"Optional authority required to authorize user"`
	TokenURL      string            `json:"token_url,omitempty"         long:"token-url"       description:"Generic OAuth provider TokenURL endpoint."`
}

func (config *GenericOAuthConfig) AuthMethod(oauthBaseURL string, teamName string) atc.AuthMethod {
	path, err := routes.OAuthRoutes.CreatePathForRoute(
		routes.OAuthBegin,
		rata.Params{"provider": ProviderName},
	)
	if err != nil {
		panic("failed to construct oauth begin handler route: " + err.Error())
	}

	path = path + fmt.Sprintf("?team_name=%s", teamName)

	return atc.AuthMethod{
		Type:        atc.AuthTypeOAuth,
		DisplayName: config.DisplayName,
		AuthURL:     oauthBaseURL + path,
	}
}

func (config *GenericOAuthConfig) IsConfigured() bool {
	return config.AuthURL != "" ||
		config.TokenURL != "" ||
		config.ClientID != "" ||
		config.ClientSecret != "" ||
		config.DisplayName != ""
}

func (config *GenericOAuthConfig) Validate() error {
	var errs *multierror.Error
	if config.ClientID == "" || config.ClientSecret == "" {
		errs = multierror.Append(
			errs,
			errors.New("must specify --generic-oauth-client-id and --generic-oauth-client-secret to use Generic OAuth."),
		)
	}
	if config.AuthURL == "" || config.TokenURL == "" {
		errs = multierror.Append(
			errs,
			errors.New("must specify --generic-oauth-auth-url and --generic-oauth-token-url to use Generic OAuth."),
		)
	}
	if config.DisplayName == "" {
		errs = multierror.Append(
			errs,
			errors.New("must specify --generic-oauth-display-name to use Generic OAuth."),
		)
	}
	return errs.ErrorOrNil()
}

type GenericTeamProvider struct{}

func (GenericTeamProvider) AddAuthGroup(group *flags.Group) provider.AuthConfig {
	flags := &GenericOAuthConfig{}

	goGroup, err := group.AddGroup("Generic OAuth Authentication (allows access to ALL authenticated users)", "", flags)
	if err != nil {
		panic(err)
	}

	goGroup.Namespace = "generic-oauth"

	return flags
}

func (GenericTeamProvider) UnmarshalConfig(config *json.RawMessage) (provider.AuthConfig, error) {
	flags := &GenericOAuthConfig{}
	if config != nil {
		err := json.Unmarshal(*config, &flags)
		if err != nil {
			return nil, err
		}
	}
	return flags, nil
}

func (GenericTeamProvider) ProviderConstructor(
	config provider.AuthConfig,
	redirectURL string,
) (provider.Provider, bool) {
	genericOAuth := config.(*GenericOAuthConfig)

	endpoint := oauth2.Endpoint{}
	if genericOAuth.AuthURL != "" && genericOAuth.TokenURL != "" {
		endpoint.AuthURL = genericOAuth.AuthURL
		endpoint.TokenURL = genericOAuth.TokenURL
	}

	var oauthVerifier verifier.Verifier
	if genericOAuth.Scope != "" && genericOAuth.Authority = "" {
		oauthVerifier = NewScopeVerifier(genericOAuth.Scope)
	} else if genericOAuth.Scope = "" && genericOAuth.Authority != "" {
		oauthVerifier = NewAuthorityVerifier(genericOAuth.Authority)
	}
	else {
		oauthVerifier = NoopVerifier{}
	}

	return Provider{
		Verifier: oauthVerifier,
		Config: ConfigOverride{
			Config: oauth2.Config{
				ClientID:     genericOAuth.ClientID,
				ClientSecret: genericOAuth.ClientSecret,
				Endpoint:     endpoint,
				RedirectURL:  redirectURL,
			},
			AuthURLParams: genericOAuth.AuthURLParams,
		},
	}, true
}

func (v NoopVerifier) Verify(logger lager.Logger, client *http.Client) (bool, error) {
	return true, nil
}

func (provider Provider) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	for key, value := range provider.Config.AuthURLParams {
		opts = append(opts, oauth2.SetAuthURLParam(key, value))

	}
	return provider.Config.AuthCodeURL(state, opts...)
}

func (provider Provider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return provider.Config.Exchange(ctx, code)
}

func (provider Provider) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	return provider.Config.Client(ctx, t)
}

func (Provider) PreTokenClient() (*http.Client, error) {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}, nil
}
