package api

import (
	"net/http"
	"path/filepath"

	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/authserver"
	"github.com/concourse/atc/api/buildserver"
	"github.com/concourse/atc/api/cliserver"
	"github.com/concourse/atc/api/configserver"
	"github.com/concourse/atc/api/containerserver"
	"github.com/concourse/atc/api/infoserver"
	"github.com/concourse/atc/api/jobserver"
	"github.com/concourse/atc/api/loglevelserver"
	"github.com/concourse/atc/api/pipelineserver"
	"github.com/concourse/atc/api/pipes"
	"github.com/concourse/atc/api/resourceserver"
	"github.com/concourse/atc/api/resourceserver/versionserver"
	"github.com/concourse/atc/api/teamserver"
	"github.com/concourse/atc/api/volumeserver"
	"github.com/concourse/atc/api/workerserver"
	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/engine"
	"github.com/concourse/atc/pipelines"
	"github.com/concourse/atc/worker"
	"github.com/concourse/atc/wrappa"
)

func NewHandler(
	logger lager.Logger,

	externalURL string,

	wrapper wrappa.Wrappa,

	tokenGenerator auth.TokenGenerator,
	providerFactory auth.ProviderFactory,
	oAuthBaseURL string,

	pipelineDBFactory db.PipelineDBFactory,
	configDB db.ConfigDB,

	authDB authserver.AuthDB,
	buildsDB buildserver.BuildsDB,
	workerDB workerserver.WorkerDB,
	containerDB containerserver.ContainerDB,
	volumesDB volumeserver.VolumesDB,
	pipeDB pipes.PipeDB,
	pipelinesDB db.PipelinesDB,
	teamDB teamserver.TeamDB,

	configValidator configserver.ConfigValidator,
	peerURL string,
	eventHandlerFactory buildserver.EventHandlerFactory,
	drain <-chan struct{},

	engine engine.Engine,
	workerClient worker.Client,

	schedulerFactory jobserver.SchedulerFactory,
	scannerFactory resourceserver.ScannerFactory,

	sink *lager.ReconfigurableSink,

	cliDownloadsDir string,
	version string,
) (http.Handler, error) {
	absCLIDownloadsDir, err := filepath.Abs(cliDownloadsDir)
	if err != nil {
		return nil, err
	}

	pipelineHandlerFactory := pipelines.NewHandlerFactory(pipelineDBFactory)

	authServer := authserver.NewServer(
		logger,
		externalURL,
		oAuthBaseURL,
		tokenGenerator,
		providerFactory,
		authDB,
	)

	buildServer := buildserver.NewServer(
		logger,
		externalURL,
		engine,
		workerClient,
		buildsDB,
		configDB,
		eventHandlerFactory,
		drain,
	)

	jobServer := jobserver.NewServer(logger, schedulerFactory, externalURL)
	resourceServer := resourceserver.NewServer(logger, scannerFactory)
	versionServer := versionserver.NewServer(logger, externalURL)
	pipeServer := pipes.NewServer(logger, peerURL, externalURL, pipeDB)

	pipelineServer := pipelineserver.NewServer(logger, pipelinesDB, configDB)

	configServer := configserver.NewServer(logger, configDB, configValidator)

	workerServer := workerserver.NewServer(logger, workerDB)

	logLevelServer := loglevelserver.NewServer(logger, sink)

	cliServer := cliserver.NewServer(logger, absCLIDownloadsDir)

	containerServer := containerserver.NewServer(logger, workerClient, containerDB)

	volumesServer := volumeserver.NewServer(logger, volumesDB)

	teamServer := teamserver.NewServer(logger, teamDB)

	infoServer := infoserver.NewServer(logger, version)

	handlers := map[string]http.Handler{
		atc.ListAuthMethods: http.HandlerFunc(authServer.ListAuthMethods),
		atc.GetAuthToken:    http.HandlerFunc(authServer.GetAuthToken),

		atc.GetConfig:  http.HandlerFunc(configServer.GetConfig),
		atc.SaveConfig: http.HandlerFunc(configServer.SaveConfig),

		atc.GetBuild:            http.HandlerFunc(buildServer.GetBuild),
		atc.ListBuilds:          http.HandlerFunc(buildServer.ListBuilds),
		atc.CreateBuild:         http.HandlerFunc(buildServer.CreateBuild),
		atc.BuildEvents:         http.HandlerFunc(buildServer.BuildEvents),
		atc.BuildResources:      http.HandlerFunc(buildServer.BuildResources),
		atc.AbortBuild:          http.HandlerFunc(buildServer.AbortBuild),
		atc.GetBuildPlan:        http.HandlerFunc(buildServer.GetBuildPlan),
		atc.GetBuildPreparation: http.HandlerFunc(buildServer.GetBuildPreparation),

		atc.ListJobs:       pipelineHandlerFactory.HandlerFor(jobServer.ListJobs),
		atc.GetJob:         pipelineHandlerFactory.HandlerFor(jobServer.GetJob),
		atc.ListJobBuilds:  pipelineHandlerFactory.HandlerFor(jobServer.ListJobBuilds),
		atc.ListJobInputs:  pipelineHandlerFactory.HandlerFor(jobServer.ListJobInputs),
		atc.GetJobBuild:    pipelineHandlerFactory.HandlerFor(jobServer.GetJobBuild),
		atc.CreateJobBuild: pipelineHandlerFactory.HandlerFor(jobServer.CreateJobBuild),
		atc.PauseJob:       pipelineHandlerFactory.HandlerFor(jobServer.PauseJob),
		atc.UnpauseJob:     pipelineHandlerFactory.HandlerFor(jobServer.UnpauseJob),
		atc.JobBadge:       pipelineHandlerFactory.HandlerFor(jobServer.JobBadge),

		atc.ListPipelines:   http.HandlerFunc(pipelineServer.ListPipelines),
		atc.GetPipeline:     http.HandlerFunc(pipelineServer.GetPipeline),
		atc.DeletePipeline:  pipelineHandlerFactory.HandlerFor(pipelineServer.DeletePipeline),
		atc.OrderPipelines:  http.HandlerFunc(pipelineServer.OrderPipelines),
		atc.PausePipeline:   pipelineHandlerFactory.HandlerFor(pipelineServer.PausePipeline),
		atc.UnpausePipeline: pipelineHandlerFactory.HandlerFor(pipelineServer.UnpausePipeline),
		atc.GetVersionsDB:   pipelineHandlerFactory.HandlerFor(pipelineServer.GetVersionsDB),
		atc.RenamePipeline:  pipelineHandlerFactory.HandlerFor(pipelineServer.RenamePipeline),

		atc.ListResources:   pipelineHandlerFactory.HandlerFor(resourceServer.ListResources),
		atc.GetResource:     pipelineHandlerFactory.HandlerFor(resourceServer.GetResource),
		atc.PauseResource:   pipelineHandlerFactory.HandlerFor(resourceServer.PauseResource),
		atc.UnpauseResource: pipelineHandlerFactory.HandlerFor(resourceServer.UnpauseResource),
		atc.CheckResource:   pipelineHandlerFactory.HandlerFor(resourceServer.CheckResource),

		atc.ListResourceVersions:          pipelineHandlerFactory.HandlerFor(versionServer.ListResourceVersions),
		atc.EnableResourceVersion:         pipelineHandlerFactory.HandlerFor(versionServer.EnableResourceVersion),
		atc.DisableResourceVersion:        pipelineHandlerFactory.HandlerFor(versionServer.DisableResourceVersion),
		atc.ListBuildsWithVersionAsInput:  pipelineHandlerFactory.HandlerFor(versionServer.ListBuildsWithVersionAsInput),
		atc.ListBuildsWithVersionAsOutput: pipelineHandlerFactory.HandlerFor(versionServer.ListBuildsWithVersionAsOutput),

		atc.CreatePipe: http.HandlerFunc(pipeServer.CreatePipe),
		atc.WritePipe:  http.HandlerFunc(pipeServer.WritePipe),
		atc.ReadPipe:   http.HandlerFunc(pipeServer.ReadPipe),

		atc.ListWorkers:    http.HandlerFunc(workerServer.ListWorkers),
		atc.RegisterWorker: http.HandlerFunc(workerServer.RegisterWorker),

		atc.SetLogLevel: http.HandlerFunc(logLevelServer.SetMinLevel),
		atc.GetLogLevel: http.HandlerFunc(logLevelServer.GetMinLevel),

		atc.DownloadCLI: http.HandlerFunc(cliServer.Download),
		atc.GetInfo:     http.HandlerFunc(infoServer.Info),

		atc.ListContainers:  http.HandlerFunc(containerServer.ListContainers),
		atc.GetContainer:    http.HandlerFunc(containerServer.GetContainer),
		atc.HijackContainer: http.HandlerFunc(containerServer.HijackContainer),

		atc.ListVolumes: http.HandlerFunc(volumesServer.ListVolumes),

		atc.SetTeam: http.HandlerFunc(teamServer.SetTeam),
	}

	return rata.NewRouter(atc.Routes, wrapper.Wrap(handlers))
}
