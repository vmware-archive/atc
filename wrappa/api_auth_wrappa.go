package wrappa

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/auth"
	"github.com/tedsuo/rata"
)

type APIAuthWrappa struct {
	AuthValidator     auth.Validator
	TokenValidator    auth.Validator
	UserContextReader auth.UserContextReader
}

func NewAPIAuthWrappa(
	authValidator auth.Validator,
	tokenValidator auth.Validator,
	userContextReader auth.UserContextReader,
) *APIAuthWrappa {
	return &APIAuthWrappa{
		AuthValidator:     authValidator,
		TokenValidator:    tokenValidator,
		UserContextReader: userContextReader,
	}
}

func (wrappa *APIAuthWrappa) Wrap(handlers rata.Handlers) rata.Handlers {
	wrapped := rata.Handlers{}

	rejector := auth.UnauthorizedRejector{}

	for name, handler := range handlers {
		newHandler := handler

		switch name {
		// unauthenticated / delegating to handler
		case atc.DownloadCLI,
			atc.ListAuthMethods,
			atc.GetInfo,
			atc.BuildEvents,
			atc.GetBuild,
			atc.BuildResources,
			atc.GetBuildPlan,
			atc.GetBuildPreparation,
			atc.ListAllPipelines,
			atc.ListBuilds,
			atc.GetJobBuild,
			atc.JobBadge,
			atc.ListJobs,
			atc.GetJob,
			atc.ListJobBuilds,
			atc.GetResource,
			atc.ListBuildsWithVersionAsInput,
			atc.ListBuildsWithVersionAsOutput,
			atc.ListResources,
			atc.ListResourceVersions,
			atc.ListPipelines,
			atc.GetPipeline,
			atc.ListTeams:

		// authenticated
		case atc.GetAuthToken,
			atc.AbortBuild,
			atc.CreateBuild,
			atc.CreatePipe,
			atc.GetContainer,
			atc.HijackContainer,
			atc.ListContainers,
			atc.ListWorkers,
			atc.ReadPipe,
			atc.RegisterWorker,
			atc.SetLogLevel,
			atc.SetTeam,
			atc.WritePipe,
			atc.ListVolumes,
			atc.GetLogLevel:
			newHandler = auth.CheckAuthenticationHandler(handler, rejector)

		// authorized
		case atc.CheckResource,
			atc.CreateJobBuild,
			atc.DeletePipeline,
			atc.DisableResourceVersion,
			atc.EnableResourceVersion,
			atc.GetConfig,
			atc.GetVersionsDB,
			atc.ListJobInputs,
			atc.OrderPipelines,
			atc.PauseJob,
			atc.PausePipeline,
			atc.PauseResource,
			atc.RenamePipeline,
			atc.UnpauseJob,
			atc.UnpausePipeline,
			atc.UnpauseResource,
			atc.RevealPipeline,
			atc.ConcealPipeline,
			atc.SaveConfig:
			newHandler = auth.CheckAuthorizationHandler(handler, rejector)

		// think about it!
		default:
			panic("you missed a spot")
		}

		if name == atc.GetAuthToken {
			newHandler = auth.WrapHandler(newHandler, wrappa.AuthValidator, wrappa.UserContextReader)
		} else {
			newHandler = auth.WrapHandler(newHandler, wrappa.TokenValidator, wrappa.UserContextReader)
		}
		wrapped[name] = newHandler
	}

	return wrapped
}
