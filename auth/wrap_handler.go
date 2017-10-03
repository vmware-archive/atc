package auth

import (
	"context"
	"net/http"

	"code.cloudfoundry.org/lager"
)

const authenticated = "authenticated"
const teamNameKey = "teamName"
const isAdminKey = "isAdmin"
const isSystemKey = "system"
const CSRFTokenKey = "csrfToken"

func WrapHandler(
	logger lager.Logger,
	handler http.Handler,
	validator Validator,
	userContextReader UserContextReader,
) http.Handler {
	return authHandler{
		logger:            logger,
		handler:           handler,
		validator:         validator,
		userContextReader: userContextReader,
	}
}

type authHandler struct {
	logger            lager.Logger
	handler           http.Handler
	validator         Validator
	userContextReader UserContextReader
}

func (h authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), authenticated, h.validator.IsAuthenticated(h.logger, r))
	teamName, isAdmin, found := h.userContextReader.GetTeam(r)
	if found {
		ctx = context.WithValue(ctx, teamNameKey, teamName)
		ctx = context.WithValue(ctx, isAdminKey, isAdmin)
	}

	isSystem, found := h.userContextReader.GetSystem(r)
	if found {
		ctx = context.WithValue(ctx, isSystemKey, isSystem)
	}
	h.handler.ServeHTTP(w, r.WithContext(ctx))
}
