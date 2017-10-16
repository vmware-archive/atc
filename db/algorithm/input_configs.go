package algorithm

type InputConfigs []InputConfig

type InputConfig struct {
	Name    string
	JobName string

	// XXX: this is *job permutation*, not *job*, since there's a subset of the
	// permutations that is implied by not mentioning other resources and their
	// spaces (the fan-in)
	Passed    JobPermutationSet
	PassedAll JobPermutationSet

	UseEveryVersion  bool
	PinnedVersionID  int
	ResourceSpaceID  int
	JobPermutationID int
}

func (configs InputConfigs) Resolve(db *VersionsDB) (InputMapping, bool) {
	jobs := JobPermutationSet{}
	inputCandidates := InputCandidates{}

	for _, inputConfig := range configs {
		versionCandidates := VersionCandidates{}

		if len(inputConfig.Passed) == 0 {
			if inputConfig.UseEveryVersion {
				versionCandidates = db.AllVersionsOfResource(inputConfig.ResourceSpaceID)
			} else {
				var versionCandidate VersionCandidate
				var found bool

				if inputConfig.PinnedVersionID != 0 {
					versionCandidate, found = db.FindVersionOfResource(inputConfig.ResourceSpaceID, inputConfig.PinnedVersionID)
				} else {
					versionCandidate, found = db.LatestVersionOfResource(inputConfig.ResourceSpaceID)
				}

				if found {
					versionCandidates.Add(versionCandidate)
				}
			}

			if versionCandidates.IsEmpty() {
				return nil, false
			}
		} else {
			jobs = jobs.Union(inputConfig.Passed)

			versionCandidates = db.VersionsOfResourcePassedJobs(
				inputConfig.ResourceSpaceID,
				inputConfig.Passed,
				inputConfig.PassedAll,
			)

			if versionCandidates.IsEmpty() {
				return nil, false
			}
		}

		existingBuildResolver := &ExistingBuildResolver{
			BuildInputs:      db.BuildInputs,
			JobPermutationID: inputConfig.JobPermutationID,
			ResourceSpaceID:  inputConfig.ResourceSpaceID,
		}

		inputCandidates = append(inputCandidates, InputVersionCandidates{
			Input:                 inputConfig.Name,
			Passed:                inputConfig.Passed,
			UseEveryVersion:       inputConfig.UseEveryVersion,
			PinnedVersionID:       inputConfig.PinnedVersionID,
			VersionCandidates:     versionCandidates,
			ExistingBuildResolver: existingBuildResolver,
		})
	}

	basicMapping, ok := inputCandidates.Reduce(0, jobs)
	if !ok {
		return nil, false
	}

	mapping := InputMapping{}
	for _, inputConfig := range configs {
		inputName := inputConfig.Name
		inputVersionID := basicMapping[inputName]
		firstOccurrence := db.IsVersionFirstOccurrence(inputVersionID, inputConfig.JobPermutationID, inputName)
		mapping[inputName] = InputVersion{
			VersionID:       inputVersionID,
			FirstOccurrence: firstOccurrence,
		}
	}

	return mapping, true
}
