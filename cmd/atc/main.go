package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/migration"
	gclient "github.com/cloudfoundry-incubator/garden/client"
	gconn "github.com/cloudfoundry-incubator/garden/client/connection"
	httpmetrics "github.com/codahale/http-handlers/metrics"
	_ "github.com/codahale/metrics/runtime"
	"github.com/lib/pq"
	"github.com/nu7hatch/gouuid"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/ifrit/sigmon"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api"
	"github.com/concourse/atc/api/buildserver"
	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/builds"
	"github.com/concourse/atc/config"
	Db "github.com/concourse/atc/db"
	"github.com/concourse/atc/db/migrations"
	"github.com/concourse/atc/engine"
	"github.com/concourse/atc/exec"
	"github.com/concourse/atc/pipelines"
	rdr "github.com/concourse/atc/radar"
	"github.com/concourse/atc/resource"
	sched "github.com/concourse/atc/scheduler"
	"github.com/concourse/atc/web"
	"github.com/concourse/atc/worker"
)

var pipelinePath = flag.String(
	"pipeline",
	"",
	"path to atc pipeline config .yml",
)

var templatesDir = flag.String(
	"templates",
	"./web/templates",
	"path to directory containing the html templates",
)

var publicDir = flag.String(
	"public",
	"./web/public",
	"path to directory containing public resources (javascript, css, etc.)",
)

var gardenNetwork = flag.String(
	"gardenNetwork",
	"",
	"garden API network type (tcp/unix). leave empty for dynamic registration.",
)

var gardenAddr = flag.String(
	"gardenAddr",
	"",
	"garden API network address (host:port or socket path). leave empty for dynamic registration.",
)

var resourceTypes = flag.String(
	"resourceTypes",
	`[
		{"type": "archive", "image": "docker:///concourse/archive-resource" },
		{"type": "docker-image", "image": "docker:///concourse/docker-image-resource" },
		{"type": "git", "image": "docker:///concourse/git-resource" },
		{"type": "github-release", "image": "docker:///concourse/github-release-resource" },
		{"type": "s3", "image": "docker:///concourse/s3-resource" },
		{"type": "semver", "image": "docker:///concourse/semver-resource" },
		{"type": "time", "image": "docker:///concourse/time-resource" },
		{"type": "tracker", "image": "docker:///concourse/tracker-resource" }
	]`,
	"map of resource type to its rootfs",
)

var sqlDriver = flag.String(
	"sqlDriver",
	"postgres",
	"database/sql driver name",
)

var sqlDataSource = flag.String(
	"sqlDataSource",
	"postgres://127.0.0.1:5432/atc?sslmode=disable",
	"database/sql data source configuration string",
)

var webListenAddress = flag.String(
	"webListenAddress",
	"0.0.0.0",
	"address to listen on",
)

var webListenPort = flag.Int(
	"webListenPort",
	8080,
	"port for the web server to listen on",
)

var callbacksURLString = flag.String(
	"callbacksURL",
	"http://127.0.0.1:8080",
	"URL used for callbacks to reach the ATC (excluding basic auth)",
)

var debugListenAddress = flag.String(
	"debugListenAddress",
	"127.0.0.1",
	"address for the pprof debugger listen on",
)

var debugListenPort = flag.Int(
	"debugListenPort",
	8079,
	"port for the pprof debugger to listen on",
)

var httpUsername = flag.String(
	"httpUsername",
	"",
	"basic auth username for the server",
)

var httpPassword = flag.String(
	"httpPassword",
	"",
	"basic auth password for the server",
)

var httpHashedPassword = flag.String(
	"httpHashedPassword",
	"",
	"bcrypted basic auth password for the server",
)

var checkInterval = flag.Duration(
	"checkInterval",
	1*time.Minute,
	"interval on which to poll for new versions of resources",
)

var publiclyViewable = flag.Bool(
	"publiclyViewable",
	false,
	"allow viewability without authentication (destructive operations still require auth)",
)

var dev = flag.Bool(
	"dev",
	false,
	"dev mode; lax security",
)

var noop = flag.Bool(
	"noop",
	false,
	"don't trigger any builds automatically",
)

var cliDownloadsDir = flag.String(
	"cliDownloadsDir",
	"",
	"directory containing CLI binaries to serve",
)

var forceHTTPS = flag.Bool(
	"forceHTTPS",
	false,
	"require client connections to use HTTPS",
)

func main() {
	flag.Parse()

	if !*dev && (*httpUsername == "" || (*httpHashedPassword == "" && *httpPassword == "")) {
		fatal(errors.New("must specify -httpUsername and -httpPassword or -httpHashedPassword or turn on dev mode"))
	}

	if _, err := os.Stat(*templatesDir); err != nil {
		fatal(errors.New("directory specified via -templates does not exist"))
	}

	if _, err := os.Stat(*publicDir); err != nil {
		fatal(errors.New("directory specified via -public does not exist"))
	}

	logger := lager.NewLogger("atc")

	logLevel := lager.INFO
	if *dev {
		logLevel = lager.DEBUG
	}

	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), logLevel)
	logger.RegisterSink(sink)

	var err error

	var dbConn *sql.DB

	for {
		dbConn, err = migration.Open(*sqlDriver, *sqlDataSource, migrations.Migrations)
		if err != nil {
			if strings.Contains(err.Error(), " dial ") {
				logger.Error("failed-to-open-db", err)
				time.Sleep(5 * time.Second)
				continue
			}

			fatal(err)
		}

		break
	}

	listener := pq.NewListener(*sqlDataSource, time.Second, time.Minute, nil)
	bus := Db.NewNotificationsBus(listener)

	db := Db.NewSQL(logger.Session("db"), dbConn, bus)
	pipelineDBFactory := Db.NewPipelineDBFactory(logger.Session("db"), dbConn, bus, db)

	var configDB Db.ConfigDB
	configDB = Db.PlanConvertingConfigDB{db}

	var resourceTypesNG []atc.WorkerResourceType
	err = json.Unmarshal([]byte(*resourceTypes), &resourceTypesNG)
	if err != nil {
		logger.Fatal("invalid-resource-types", err)
	}

	var workerClient worker.Client
	if *gardenAddr != "" {
		workerClient = worker.NewGardenWorker(
			gclient.New(gconn.NewWithLogger(
				*gardenNetwork,
				*gardenAddr,
				logger.Session("garden-connection"),
			)),
			clock.NewClock(),
			-1,
			resourceTypesNG,
			"linux",
			[]string{},
		)
	} else {
		workerClient = worker.NewPool(worker.NewDBWorkerProvider(db, logger))
	}

	resourceTracker := resource.NewTracker(workerClient)
	gardenFactory := exec.NewGardenFactory(workerClient, resourceTracker, func() string {
		guid, err := uuid.NewV4()
		if err != nil {
			panic("not enough entropy to generate guid: " + err.Error())
		}

		return guid.String()
	})
	execEngine := engine.NewExecEngine(gardenFactory, engine.NewBuildDelegateFactory(db), db)

	engine := engine.NewDBEngine(engine.Engines{execEngine}, db, db)

	var webValidator auth.Validator

	if *httpUsername != "" && *httpHashedPassword != "" {
		webValidator = auth.BasicAuthHashedValidator{
			Username:       *httpUsername,
			HashedPassword: *httpHashedPassword,
		}
	} else if *httpUsername != "" && *httpPassword != "" {
		webValidator = auth.BasicAuthValidator{
			Username: *httpUsername,
			Password: *httpPassword,
		}
	} else {
		webValidator = auth.NoopValidator{}
	}

	callbacksURL, err := url.Parse(*callbacksURLString)
	if err != nil {
		fatal(err)
	}

	drain := make(chan struct{})

	apiHandler, err := api.NewHandler(
		logger,            // logger lager.Logger,
		webValidator,      // validator auth.Validator,
		pipelineDBFactory, // pipelineDBFactory db.PipelineDBFactory,

		configDB, // configDB db.ConfigDB,

		db, // buildsDB buildserver.BuildsDB,
		db, // workerDB workerserver.WorkerDB,
		db, // pipeDB pipes.PipeDB,
		db, // pipelinesDB db.PipelinesDB,

		config.ValidateConfig,       // configValidator configserver.ConfigValidator,
		callbacksURL.String(),       // peerURL string,
		buildserver.NewEventHandler, // eventHandlerFactory buildserver.EventHandlerFactory,
		drain, // drain <-chan struct{},

		engine,       // engine engine.Engine,
		workerClient, // workerClient worker.Client,

		sink, // sink *lager.ReconfigurableSink,

		*cliDownloadsDir, // cliDownloadsDir string,
		*forceHTTPS,
	)
	if err != nil {
		fatal(err)
	}

	radarSchedulerFactory := pipelines.NewRadarSchedulerFactory(
		resourceTracker,
		*checkInterval,
		db,
		engine,
		db,
	)

	webHandler, err := web.NewHandler(
		logger,
		webValidator,
		radarSchedulerFactory,
		db,
		pipelineDBFactory,
		configDB,
		*templatesDir,
		*publicDir,
		engine,
		*forceHTTPS,
	)
	if err != nil {
		fatal(err)
	}

	webMux := http.NewServeMux()
	webMux.Handle("/api/v1/", apiHandler)
	webMux.Handle("/", webHandler)

	var httpHandler http.Handler

	httpHandler = webMux

	if !*publiclyViewable {
		httpHandler = auth.Handler{
			Handler:   httpHandler,
			Validator: webValidator,
		}
	}

	// copy Authorization header as ATC-Authorization cookie for websocket auth
	httpHandler = auth.CookieSetHandler{
		Handler: httpHandler,
	}

	httpHandler = httpmetrics.Wrap(httpHandler)

	webListenAddr := fmt.Sprintf("%s:%d", *webListenAddress, *webListenPort)
	debugListenAddr := fmt.Sprintf("%s:%d", *debugListenAddress, *debugListenPort)

	syncer := pipelines.NewSyncer(
		logger.Session("syncer"),
		db,
		pipelineDBFactory,
		func(pipelineDB Db.PipelineDB) ifrit.Runner {
			return grouper.NewParallel(os.Interrupt, grouper.Members{
				{
					pipelineDB.ScopedName("radar"),
					rdr.NewRunner(
						logger.Session(pipelineDB.ScopedName("radar")),
						*noop,
						db,
						radarSchedulerFactory.BuildRadar(pipelineDB),
						pipelineDB,
						1*time.Minute,
					),
				},
				{
					pipelineDB.ScopedName("scheduler"),
					&sched.Runner{
						Logger: logger.Session(pipelineDB.ScopedName("scheduler")),

						Locker: db,
						DB:     pipelineDB,

						Scheduler: radarSchedulerFactory.BuildScheduler(pipelineDB),

						Noop: *noop,

						Interval: 10 * time.Second,
					},
				},
			})
		},
	)

	buildTracker := builds.NewTracker(
		logger.Session("build-tracker"),
		db,
		engine,
	)

	memberGrouper := []grouper.Member{
		{"web", http_server.New(webListenAddr, httpHandler)},

		{"debug", http_server.New(debugListenAddr, http.DefaultServeMux)},

		{"drainer", ifrit.RunFunc(func(signals <-chan os.Signal, ready chan<- struct{}) error {
			close(ready)

			<-signals

			close(drain)

			return nil
		})},

		{"pipelines", pipelines.SyncRunner{
			Syncer:   syncer,
			Interval: 10 * time.Second,
			Clock:    clock.NewClock(),
		}},

		{"builds", builds.TrackerRunner{
			Tracker:  buildTracker,
			Interval: 10 * time.Second,
			Clock:    clock.NewClock(),
		}},
	}

	group := grouper.NewParallel(os.Interrupt, memberGrouper)

	running := ifrit.Envoke(sigmon.New(group))

	logger.Info("listening", lager.Data{
		"web":   webListenAddr,
		"debug": debugListenAddr,
	})

	err = <-running.Wait()
	if err != nil {
		logger.Error("exited-with-failure", err)
		os.Exit(1)
	}
}

func fatal(err error) {
	println(err.Error())
	os.Exit(1)
}
