package wrappa

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/api/auth"
	"github.com/tedsuo/rata"
)

type APIAuthWrappa struct {
	authValidator                       auth.Validator
	userContextReader                   auth.UserContextReader
	checkPipelineAccessHandlerFactory   auth.CheckPipelineAccessHandlerFactory
	checkBuildReadAccessHandlerFactory  auth.CheckBuildReadAccessHandlerFactory
	checkBuildWriteAccessHandlerFactory auth.CheckBuildWriteAccessHandlerFactory
	checkWorkerTeamAccessHandlerFactory auth.CheckWorkerTeamAccessHandlerFactory
}

func NewAPIAuthWrappa(
	authValidator auth.Validator,
	userContextReader auth.UserContextReader,
	checkPipelineAccessHandlerFactory auth.CheckPipelineAccessHandlerFactory,
	checkBuildReadAccessHandlerFactory auth.CheckBuildReadAccessHandlerFactory,
	checkBuildWriteAccessHandlerFactory auth.CheckBuildWriteAccessHandlerFactory,
	checkWorkerTeamAccessHandlerFactory auth.CheckWorkerTeamAccessHandlerFactory,
) *APIAuthWrappa {
	return &APIAuthWrappa{
		authValidator:                       authValidator,
		userContextReader:                   userContextReader,
		checkPipelineAccessHandlerFactory:   checkPipelineAccessHandlerFactory,
		checkBuildReadAccessHandlerFactory:  checkBuildReadAccessHandlerFactory,
		checkBuildWriteAccessHandlerFactory: checkBuildWriteAccessHandlerFactory,
		checkWorkerTeamAccessHandlerFactory: checkWorkerTeamAccessHandlerFactory,
	}
}

func (wrappa *APIAuthWrappa) Wrap(handlers rata.Handlers) rata.Handlers {
	wrapped := rata.Handlers{}

	rejector := auth.UnauthorizedRejector{}

	for name, handler := range handlers {
		newHandler := handler

		switch name {
		// unauthenticated / delegating to handler
		case
			atc.ListAllPipelines,
			atc.ListPipelines,
			atc.ListBuilds,
			atc.LegacyListAuthMethods,
			atc.LegacyGetAuthToken,
			atc.LegacyGetUser,
			atc.MainJobBadge:

		// !!!!!!!!!!!!!! REFACTORED !!!!!!!!!!!!!!!!!!!!!
		case
			atc.DestroyTeam,
			atc.SetTeam,
			atc.ListTeams,
			atc.GetInfo,
			atc.DownloadCLI,
			atc.CheckResourceWebHook,
			atc.CheckResource,
			atc.UnpauseResource,
			atc.PauseResource,
			atc.GetResource,
			atc.ListResources,
			atc.CreatePipelineBuild,
			atc.GetVersionsDB,
			atc.OrderPipelines,
			atc.RenamePipeline,
			atc.HidePipeline,
			atc.ExposePipeline,
			atc.UnpausePipeline,
			atc.PausePipeline,
			atc.DeletePipeline,
			atc.PipelineBadge,
			atc.GetPipeline,
			atc.JobBadge,
			atc.UnpauseJob,
			atc.PauseJob,
			atc.ListJobInputs,
			atc.ListJobBuilds,
			atc.GetJob,
			atc.ListJobs,
			atc.GetJobBuild,
			atc.CreateJobBuild,
			atc.GetConfig,
			atc.SaveConfig:

		// authorized (requested team matches resource team)
		case
			atc.DisableResourceVersion,
			atc.EnableResourceVersion:
			newHandler = auth.CheckAuthorizationHandler(handler, rejector)

		// pipeline is public or authorized
		case
			atc.ListBuildsWithVersionAsInput,
			atc.ListBuildsWithVersionAsOutput,
			atc.GetResourceCausality,
			atc.GetResourceVersion,
			atc.ListResourceVersions:
			newHandler = wrappa.checkPipelineAccessHandlerFactory.HandlerFor(handler, rejector)

		// pipeline is public or authorized
		case
			atc.GetBuild,
			atc.BuildResources,
			atc.GetBuildPlan:
			newHandler = wrappa.checkBuildReadAccessHandlerFactory.AnyJobHandler(handler, rejector)

		// pipeline and job are public or authorized
		case atc.GetBuildPreparation,
			atc.BuildEvents:
			newHandler = wrappa.checkBuildReadAccessHandlerFactory.CheckIfPrivateJobHandler(handler, rejector)

		// resource belongs to authorized team
		case atc.AbortBuild:
			newHandler = wrappa.checkBuildWriteAccessHandlerFactory.HandlerFor(handler, rejector)

		// requester is system, admin team, or worker owning team
		case atc.PruneWorker,
			atc.LandWorker,
			atc.RetireWorker:
			newHandler = wrappa.checkWorkerTeamAccessHandlerFactory.HandlerFor(handler, rejector)

		// authenticated
		case
			atc.CreateBuild,
			atc.CreatePipe,
			atc.GetContainer,
			atc.HijackContainer,
			atc.ListContainers,
			atc.ListWorkers,
			atc.ReadPipe,
			atc.RegisterWorker,
			atc.HeartbeatWorker,
			atc.DeleteWorker,
			atc.RenameTeam,
			atc.WritePipe,
			atc.ListVolumes:
			newHandler = auth.CheckAuthenticationHandler(handler, rejector)

		case atc.GetLogLevel,
			atc.SetLogLevel:
			newHandler = auth.CheckAdminHandler(handler, rejector)

		// think about it!
		default:
			panic("you missed a spot")
		}

		newHandler = auth.WrapHandler(newHandler, wrappa.authValidator, wrappa.userContextReader)
		wrapped[name] = auth.CSRFValidationHandler(newHandler, rejector, wrappa.userContextReader)
	}

	return wrapped
}
