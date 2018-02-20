package accessor

import (
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/lock"
)

type Access int

const (
	Skip  = Access(iota)
	Read  = Access(iota)
	Write = Access(iota)
)

//go:generate counterfeiter . Accessor
type Accessor interface {
	Admin(Access) (db.Admin, error)
	Pipelines(Access) ([]db.Pipeline, error)
	Teams(Access) ([]db.Team, error)
	Team(Access, string) (db.Team, error)
	TeamPipelines(Access, string) ([]db.Pipeline, error)
	TeamPipeline(Access, string, string) (db.Pipeline, error)
	TeamPipelineJobs(Access, string, string) ([]db.Job, error)
	TeamPipelineJob(Access, string, string, string) (db.Job, error)
	TeamPipelineJobBuild(Access, string, string, string, string) (db.Build, error)
	CreateTeampipelinJobBuilds(string, string, string, string)
	TeamPipelineResources(Access, string, string) ([]db.Resource, error)
	TeamPipelineResource(Access, string, string, string) (db.Resource, error)
}

type accessor struct {
	conn        db.Conn
	lockFactory lock.LockFactory

	userId    string
	userName  string
	teamNames []string
}

func (u *accessor) Admin(a Access) (db.Admin, error) {
	return nil, nil
}

func (u *accessor) Pipelines(a Access) ([]db.Pipeline, error) {
	// var pipelines []db.Pipeline
	// if authTeamFound {
	// 	team, found, err := s.teamFactory.FindTeam(authTeam.Name())
	// 	if err != nil {
	// 		logger.Error("failed-to-get-team", err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}

	// 	if !found {
	// 		logger.Info("team-not-found")
	// 		w.WriteHeader(http.StatusNotFound)
	// 		return
	// 	}

	// 	pipelines, err = team.VisiblePipelines()
	// 	if err != nil {
	// 		logger.Error("failed-to-get-all-visible-pipelines", err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// } else {
	// 	var err error
	// 	pipelines, err = s.pipelineFactory.PublicPipelines()
	// 	if err != nil {
	// 		logger.Error("failed-to-get-all-public-pipelines", err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	return nil, nil
}

func (u *accessor) Team(a Access, teamName string) (db.Team, error) {
	return nil, nil
}

func (u *accessor) Teams(a Access) ([]db.Team, error) {
	return nil, nil
}

func (u *accessor) TeamPipelines(a Access) ([]db.Pipeline, error) {

	// var pipelines []db.Pipeline

	// authTeam, authTeamFound := auth.GetTeam(r)
	// if authTeamFound && authTeam.IsAuthorized(requestTeamName) {
	// 	pipelines, err = team.Pipelines()
	// } else {
	// 	pipelines, err = team.PublicPipelines()
	// }

	// if err != nil {
	// 	logger.Error("failed-to-get-all-active-pipelines", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	return nil, nil
}

func (u *accessor) TeamPipeline(a Access, teamName, pipelineName string) (db.Pipeline, error) {
	//NOTE: For webhook token endpoint , skip the validation for user with permission SKIP
	return nil, nil
}

func (u *accessor) TeamPipelineJobs(a Access, teamName, pipelineName string) ([]db.Job, error) {
	return []db.Job{}, nil
}

func (u *accessor) TeamPipelineJob(a Access, teamName, pipelineName, jobName string) (db.Job, error) {
	// job, found, err := pipeline.Job(jobName)
	// if err != nil {
	// 	logger.Error("failed-to-get-job", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// if !found {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	return nil, nil
}

func (u *accessor) TeamPipelineJobBuild(a Access, teamName, pipelineName, jobName, buildName string) (db.Build, error) {
	// build, found, err := job.Build(buildName)
	// if err != nil {
	// 	logger.Error("failed-to-get-job-build", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// if !found {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	return nil, nil
}

func (u *accessor) TeamPipelineResource(a Access, teamName, pipelineName string) (db.Resource, error) {

	// only include the checkErrorString with the user is authorized.

	return nil, nil
}

func (u *accessor) TeamPipelineResources(a Access, teamName, pipelineName string) ([]db.Resource, error) {
	return nil, nil
}
