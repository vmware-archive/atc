package present

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

func SpaceJobCombination(combination db.SpaceJobCombination) atc.SpaceJobCombination {
	var presentedNextBuild, presentedFinishedBuild, presentedTransitionBuild *atc.Build

	if combination.NextBuild != nil {
		presented := Build(combination.NextBuild)
		presentedNextBuild = &presented
	}

	if combination.FinishedBuild != nil {
		presented := Build(combination.FinishedBuild)
		presentedFinishedBuild = &presented
	}

	if combination.TransitionBuild != nil {
		presented := Build(combination.TransitionBuild)
		presentedTransitionBuild = &presented
	}

	return atc.SpaceJobCombination{
		ID:          combination.JobCombination.ID(),
		Combination: combination.JobCombination.Combination(),

		FinishedBuild:   presentedFinishedBuild,
		NextBuild:       presentedNextBuild,
		TransitionBuild: presentedTransitionBuild,
	}
}
