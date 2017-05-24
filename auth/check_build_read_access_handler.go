package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/concourse/atc/dbng"
)

type CheckBuildReadAccessHandlerFactory interface {
	AnyJobHandler(delegateHandler http.Handler, rejector Rejector) http.Handler
	CheckIfPrivateJobHandler(delegateHandler http.Handler, rejector Rejector) http.Handler
}

type checkBuildReadAccessHandlerFactory struct {
	buildFactory dbng.BuildFactory
}

func NewCheckBuildReadAccessHandlerFactory(
	buildFactory dbng.BuildFactory,
) *checkBuildReadAccessHandlerFactory {
	return &checkBuildReadAccessHandlerFactory{
		buildFactory: buildFactory,
	}
}

func (f *checkBuildReadAccessHandlerFactory) AnyJobHandler(
	delegateHandler http.Handler,
	rejector Rejector,
) http.Handler {
	return checkBuildReadAccessHandler{
		rejector:        rejector,
		buildFactory:    f.buildFactory,
		delegateHandler: delegateHandler,
		allowPrivateJob: true,
	}
}

func (f *checkBuildReadAccessHandlerFactory) CheckIfPrivateJobHandler(
	delegateHandler http.Handler,
	rejector Rejector,
) http.Handler {
	return checkBuildReadAccessHandler{
		rejector:        rejector,
		buildFactory:    f.buildFactory,
		delegateHandler: delegateHandler,
		allowPrivateJob: false,
	}
}

type checkBuildReadAccessHandler struct {
	rejector        Rejector
	buildFactory    dbng.BuildFactory
	delegateHandler http.Handler
	allowPrivateJob bool
}

func (h checkBuildReadAccessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.FormValue(":build_id")
	buildID, err := strconv.Atoi(buildIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	build, found, err := h.buildFactory.Build(buildID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	authTeam, authTeamFound := GetTeam(r)
	if !IsAuthenticated(r) || (authTeamFound && !authTeam.IsAuthorized(build.TeamName())) {
		pipeline, found, err := build.Pipeline()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !found {
			h.rejector.Unauthorized(w, r)
			return
		}

		if !pipeline.Public() {
			if IsAuthenticated(r) {
				h.rejector.Forbidden(w, r)
				return
			}

			h.rejector.Unauthorized(w, r)
			return
		}

		if !h.allowPrivateJob {
			config, _, _, err := pipeline.Config()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			isJobPublic, err := config.JobIsPublic(build.JobName())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !isJobPublic {
				if IsAuthenticated(r) {
					h.rejector.Forbidden(w, r)
					return
				}

				h.rejector.Unauthorized(w, r)
				return
			}
		}
	}

	ctx := context.WithValue(r.Context(), BuildContextKey, build)
	h.delegateHandler.ServeHTTP(w, r.WithContext(ctx))
}
