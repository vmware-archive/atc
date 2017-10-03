package auth

import (
	"crypto/rsa"
	"net/http"

	"code.cloudfoundry.org/lager"
)

type JWTValidator struct {
	PublicKey *rsa.PublicKey
}

func (validator JWTValidator) IsAuthenticated(logger lager.Logger, r *http.Request) bool {
	token, err := getJWT(r, validator.PublicKey)
	if err != nil {
		return false
	}

	return token.Valid
}
