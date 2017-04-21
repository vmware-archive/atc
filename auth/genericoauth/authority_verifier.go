package genericoauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/oauth2"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/auth/verifier"
)

type AuthorityVerifier struct {
	authority string
}

func NewAuthorityVerifier(
	authority string,
) verifier.Verifier {
	return AuthorityVerifier{
		authority: authority,
	}
}

type GenericOAuthToken struct {
	Authorities []string `json:"authorities"`
}

func (verifier AuthorityVerifier) Verify(logger lager.Logger, httpClient *http.Client) (bool, error) {
	oauth2Transport, ok := httpClient.Transport.(*oauth2.Transport)
	if !ok {
		return false, errors.New("httpClient transport must be of type oauth2.Transport")
	}

	token, err := oauth2Transport.Source.Token()
	if err != nil {
		return false, err
	}

	tokenParts := strings.Split(token.AccessToken, ".")
	if len(tokenParts) < 2 {
		return false, errors.New("access token contains an invalid number of segments")
	}

	decodedClaims, err := jwt.DecodeSegment(tokenParts[1])
	if err != nil {
		return false, err
	}

	var oauthToken GenericOAuthToken
	err = json.Unmarshal(decodedClaims, &oauthToken)
	if err != nil {
		return false, err
	}

	if len(oauthToken.Authorities) == 0 {
		return false, errors.New("user has no assigned authorities in access token")
	}

	for _, userAuthority := range oauthToken.Authorities {
		if userAuthority == verifier.authority {
			return true, nil
		}
	}

	logger.Info("does-not-have-authority", lager.Data{
		"have": oauthToken.Authorities,
		"want": verifier.authority,
	})

	return false, nil
}
