package db

type SpaceJobCombination struct {
	JobCombination JobCombination

	FinishedBuild   Build
	NextBuild       Build
	TransitionBuild Build
}

type SpaceJob struct {
	Job Job

	SpaceJobCombinations []SpaceJobCombination
}

type SpaceDashboard []SpaceJob
