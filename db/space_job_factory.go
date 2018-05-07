package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc/db/lock"
)

//go:generate counterfeiter . JobFactory

type SpaceJobFactory interface {
	PipelineJobs(string, string) (SpaceDashboard, error)
	VisibleJobs([]string) (SpaceDashboard, error)
}

type spaceJobFactory struct {
	conn        Conn
	lockFactory lock.LockFactory
}

func NewSpaceJobFactory(conn Conn, lockFactory lock.LockFactory) SpaceJobFactory {
	return &spaceJobFactory{
		conn:        conn,
		lockFactory: lockFactory,
	}
}

func (j *spaceJobFactory) PipelineJobs(teamName string, pipelineName string) (SpaceDashboard, error) {
	rows, err := jobsQuery.
		Where(sq.Eq{
			"t.name":   teamName,
			"p.name":   pipelineName,
			"j.active": true,
		}).
		OrderBy("j.id ASC").
		RunWith(j.conn).
		Query()
	if err != nil {
		return nil, err
	}

	jobs, err := scanJobs(j.conn, j.lockFactory, rows)
	if err != nil {
		return nil, err
	}

	var jobIDs []int
	for _, job := range jobs {
		jobIDs = append(jobIDs, job.ID())
	}

	rows, err = jobCombinationsQuery.
		Where(sq.Eq{
			"c.job_id": jobIDs,
		}).
		OrderBy("c.combination ASC").
		RunWith(j.conn).
		Query()
	if err != nil {
		return nil, err
	}

	jobCombinations, err := scanJobCombinations(j.conn, j.lockFactory, rows)
	if err != nil {
		return nil, err
	}

	var jobCombinationIDs []int
	jobCombinationsMap := make(map[int]JobCombinations)
	for _, jobCombination := range jobCombinations {
		jobCombinationIDs = append(jobCombinationIDs, jobCombination.ID())
		jobCombinationsMap[jobCombination.JobID()] = append(jobCombinationsMap[jobCombination.JobID()], jobCombination)
	}

	nextBuilds, err := j.getBuildsFrom("next_builds_per_job_combination", jobCombinationIDs)
	if err != nil {
		return nil, err
	}

	finishedBuilds, err := j.getBuildsFrom("latest_completed_builds_per_job_combination", jobCombinationIDs)
	if err != nil {
		return nil, err
	}

	transitionBuilds, err := j.getBuildsFrom("transition_builds_per_job_combination", jobCombinationIDs)
	if err != nil {
		return nil, err
	}

	spaceDashboard := SpaceDashboard{}

	for _, job := range jobs {
		spaceJob := SpaceJob{Job: job}

		spaceJobCombinations := []SpaceJobCombination{}

		for _, jobCombination := range jobCombinationsMap[job.ID()] {
			spaceJobCombination := SpaceJobCombination{JobCombination: jobCombination}

			if nextBuild, found := nextBuilds[jobCombination.ID()]; found {
				spaceJobCombination.NextBuild = nextBuild
			}

			if finishedBuild, found := finishedBuilds[jobCombination.ID()]; found {
				spaceJobCombination.FinishedBuild = finishedBuild
			}

			if transitionBuild, found := transitionBuilds[jobCombination.ID()]; found {
				spaceJobCombination.TransitionBuild = transitionBuild
			}

			spaceJobCombinations = append(spaceJobCombinations, spaceJobCombination)
		}

		spaceJob.SpaceJobCombinations = spaceJobCombinations

		spaceDashboard = append(spaceDashboard, spaceJob)
	}

	return spaceDashboard, nil
}

func (j *spaceJobFactory) VisibleJobs(teamNames []string) (SpaceDashboard, error) {
	rows, err := jobsQuery.
		Where(sq.Eq{
			"t.name":   teamNames,
			"j.active": true,
		}).
		OrderBy("j.id ASC").
		RunWith(j.conn).
		Query()
	if err != nil {
		return nil, err
	}

	currentTeamJobs, err := scanJobs(j.conn, j.lockFactory, rows)
	if err != nil {
		return nil, err
	}

	rows, err = jobsQuery.
		Where(sq.NotEq{
			"t.name": teamNames,
		}).
		Where(sq.Eq{
			"p.public": true,
			"j.active": true,
		}).
		OrderBy("j.id ASC").
		RunWith(j.conn).
		Query()
	if err != nil {
		return nil, err
	}

	otherTeamPublicJobs, err := scanJobs(j.conn, j.lockFactory, rows)
	if err != nil {
		return nil, err
	}

	jobs := append(currentTeamJobs, otherTeamPublicJobs...)

	var jobIDs []int
	for _, job := range jobs {
		jobIDs = append(jobIDs, job.ID())
	}

	rows, err = jobCombinationsQuery.
		Where(sq.Eq{
			"c.job_id": jobIDs,
		}).
		RunWith(j.conn).
		Query()
	if err != nil {
		return nil, err
	}

	jobCombinations, err := scanJobCombinations(j.conn, j.lockFactory, rows)
	if err != nil {
		return nil, err
	}

	var jobCombinationIDs []int
	jobCombinationsMap := make(map[int]JobCombinations)
	for _, jobCombination := range jobCombinations {
		jobCombinationIDs = append(jobCombinationIDs, jobCombination.ID())
		jobCombinationsMap[jobCombination.JobID()] = append(jobCombinationsMap[jobCombination.JobID()], jobCombination)
	}

	nextBuilds, err := j.getBuildsFrom("next_builds_per_job_combination", jobCombinationIDs)
	if err != nil {
		return nil, err
	}

	finishedBuilds, err := j.getBuildsFrom("latest_completed_builds_per_job_combination", jobCombinationIDs)
	if err != nil {
		return nil, err
	}

	transitionBuilds, err := j.getBuildsFrom("transition_builds_per_job_combination", jobCombinationIDs)
	if err != nil {
		return nil, err
	}

	spaceDashboard := SpaceDashboard{}

	for _, job := range jobs {
		spaceJob := SpaceJob{Job: job}

		spaceJobCombinations := []SpaceJobCombination{}

		for _, jobCombination := range jobCombinationsMap[job.ID()] {
			spaceJobCombination := SpaceJobCombination{JobCombination: jobCombination}

			if nextBuild, found := nextBuilds[jobCombination.ID()]; found {
				spaceJobCombination.NextBuild = nextBuild
			}

			if finishedBuild, found := finishedBuilds[jobCombination.ID()]; found {
				spaceJobCombination.FinishedBuild = finishedBuild
			}

			if transitionBuild, found := transitionBuilds[jobCombination.ID()]; found {
				spaceJobCombination.TransitionBuild = transitionBuild
			}

			spaceJobCombinations = append(spaceJobCombinations, spaceJobCombination)
		}

		spaceJob.SpaceJobCombinations = spaceJobCombinations

		spaceDashboard = append(spaceDashboard, spaceJob)
	}

	return spaceDashboard, nil
}

func (j *spaceJobFactory) getBuildsFrom(view string, jobCombinationIDs []int) (map[int]Build, error) {
	rows, err := buildsQuery.
		From(view + " b").
		Where(sq.Eq{"c.id": jobCombinationIDs}).
		RunWith(j.conn).Query()
	if err != nil {
		return nil, err
	}

	defer Close(rows)

	builds := make(map[int]Build)

	for rows.Next() {
		build := &build{conn: j.conn, lockFactory: j.lockFactory}
		err := scanBuild(build, rows, j.conn.EncryptionStrategy())
		if err != nil {
			return nil, err
		}
		builds[build.JobCombinationID()] = build
	}

	return builds, nil
}
