package present

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/web"
	"github.com/tedsuo/rata"
)

func Pipeline(savedPipeline db.SavedPipeline, config atc.Config) atc.Pipeline {
	pathForRoute, err := web.Routes.CreatePathForRoute(web.Pipeline, rata.Params{
		"team_name": savedPipeline.TeamName,
		"pipeline":  savedPipeline.Name,
	})

	if err != nil {
		panic("failed to generate url: " + err.Error())
	}

	return atc.Pipeline{
		Name:     savedPipeline.Name,
		TeamName: savedPipeline.TeamName,
		URL:      pathForRoute,
		Paused:   savedPipeline.Paused,
		Public:   savedPipeline.Public,
		Groups:   config.Groups,
	}
}

func Pipelines(savedPipelines []db.SavedPipeline) []atc.Pipeline {
	pipelines := make([]atc.Pipeline, len(savedPipelines))

	for i := range savedPipelines {
		pipelines[i] = Pipeline(savedPipelines[i], savedPipelines[i].Config)
	}

	return pipelines
}
