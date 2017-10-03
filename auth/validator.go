package auth

import (
	"net/http"

	"code.cloudfoundry.org/lager"
)

//go:generate counterfeiter . Validator
type Validator interface {
	IsAuthenticated(lager.Logger, *http.Request) bool
}
