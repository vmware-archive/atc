package buildserver

import (
	"net/http"
	"strconv"

	"github.com/concourse/atc/auth"
)

func (s *Server) BuildEvents(w http.ResponseWriter, r *http.Request) {
	buildID, err := strconv.Atoi(r.FormValue(":build_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	build, err := s.db.GetBuild(buildID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !s.fallback.IsAuthenticated(r) {
		if build.OneOff() {
			auth.Unauthorized(w)
			return
		}

		public, err := s.db.JobIsPublic(build.JobName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !public {
			auth.Unauthorized(w)
			return
		}
	}

	streamDone := make(chan struct{})

	go func() {
		defer close(streamDone)
		s.eventHandlerFactory(s.db, buildID, nil).ServeHTTP(w, r)
	}()

	select {
	case <-streamDone:
	case <-s.drain:
	}
}
