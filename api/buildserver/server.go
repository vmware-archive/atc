package buildserver

import (
	"net/http"
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/builder"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/event"
	"github.com/pivotal-golang/lager"
)

type EventHandlerFactory func(event.BuildsDB, int, event.Censor) http.Handler

type server struct {
	logger lager.Logger

	db                  BuildsDB
	configDB            ConfigDB
	builder             builder.Builder
	pingInterval        time.Duration
	eventHandlerFactory EventHandlerFactory
	drain               <-chan struct{}
	fallback            auth.Validator

	httpClient *http.Client
}

type BuildsDB interface {
	GetBuild(buildID int) (db.Build, error)
	GetAllBuilds() ([]db.Build, error)

	CreateOneOffBuild() (db.Build, error)
	SaveBuildStatus(buildID int, status db.Status) error

	GetBuildEvents(buildID int) ([]db.BuildEvent, error)
}

type ConfigDB interface {
	GetConfig() (atc.Config, error)
}

func NewServer(
	logger lager.Logger,
	db BuildsDB,
	configDB ConfigDB,
	builder builder.Builder,
	pingInterval time.Duration,
	eventHandlerFactory EventHandlerFactory,
	drain <-chan struct{},
	fallback auth.Validator,
) *server {
	return &server{
		logger:              logger,
		db:                  db,
		configDB:            configDB,
		builder:             builder,
		pingInterval:        pingInterval,
		eventHandlerFactory: eventHandlerFactory,
		fallback:            fallback,

		httpClient: &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: 5 * time.Minute,
			},
		},
	}
}
