package jobserver

import (
	"fmt"
	"net/http"

	"github.com/concourse/atc/db"
)

var (
	badgePassing = badge{width: 88, fillColor: `#44cc11`, status: `passing`}
	badgeFailing = badge{width: 80, fillColor: `#e05d44`, status: `failing`}
	badgeUnknown = badge{width: 98, fillColor: `#9f9f9f`, status: `unknown`}
	badgeAborted = badge{width: 90, fillColor: `#8f4b2d`, status: `aborted`}
	badgeErrored = badge{width: 88, fillColor: `#fe7d37`, status: `errored`}
)

type badge struct {
	width     int
	fillColor string
	status    string
}

func (b *badge) statusWidth() int {
	return b.width - 37
}

func (b *badge) statusTextWidth() float64 {
	return float64(b.width)/2 + 17.5
}

func (b *badge) getSvg() string {
	const svgTemplate = `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="20"><linearGradient id="b" x2="0" y2="100%%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><mask id="a"><rect width="%d" height="20" rx="3" fill="#fff"/></mask><g mask="url(#a)"><path fill="#555" d="M0 0h37v20H0z"/><path fill="%s" d="M37 0h%dv20H37z"/><path fill="url(#b)" d="M0 0h%dv20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="18.5" y="15" fill="#010101" fill-opacity=".3">build</text><text x="18.5" y="14">build</text><text x="%.1f" y="15" fill="#010101" fill-opacity=".3">%s</text><text x="%.1f" y="14">%s</text></g></svg>`
	return fmt.Sprintf(svgTemplate, b.width, b.width, b.fillColor, b.statusWidth(), b.width, b.statusTextWidth(), b.status, b.statusTextWidth(), b.status)
}

func statusSvg(finished *db.Build) string {
	switch {
	case finished == nil:
		return badgeUnknown.getSvg()
	case finished.Status == db.StatusSucceeded:
		return badgePassing.getSvg()
	case finished.Status == db.StatusFailed:
		return badgeFailing.getSvg()
	case finished.Status == db.StatusAborted:
		return badgeAborted.getSvg()
	case finished.Status == db.StatusErrored:
		return badgeErrored.getSvg()
	}
	return badgeUnknown.getSvg()
}

func (s *Server) JobBadge(pipelineDB db.PipelineDB) http.Handler {
	logger := s.logger.Session("job-badge")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobName := r.FormValue(":job_name")

		config, _, found, err := pipelineDB.GetConfig()
		if err != nil {
			logger.Error("could-not-get-pipeline-config", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, found = config.Jobs.Lookup(jobName)
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		finished, _, err := pipelineDB.GetJobFinishedAndNextBuild(jobName)
		if err != nil {
			logger.Error("could-not-get-job-finished-and-next-build", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "image/svg+xml")

		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, statusSvg(finished))
	})
}
