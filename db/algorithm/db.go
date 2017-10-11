package algorithm

type VersionsDB struct {
	ResourceVersions []ResourceVersion

	BuildOutputs []BuildOutput
	BuildInputs  []BuildInput

	// these are only present so that we can export the data set and write tests
	// against it; it's not used for anything in production code
	ResourceSpaceIDs  map[string]int
	JobPermutationIDs map[string]int
}

type ResourceVersion struct {
	VersionID       int
	ResourceSpaceID int
	CheckOrder      int
}

type BuildOutput struct {
	ResourceVersion
	BuildID          int
	JobPermutationID int
}

type BuildInput struct {
	ResourceVersion
	BuildID          int
	JobPermutationID int
	InputName        string
}

func (db VersionsDB) IsVersionFirstOccurrence(versionID int, jobID int, inputName string) bool {
	for _, buildInput := range db.BuildInputs {
		if buildInput.VersionID == versionID &&
			buildInput.JobPermutationID == jobID &&
			buildInput.InputName == inputName {
			return false
		}
	}
	return true
}

func (db VersionsDB) AllVersionsOfResource(resourceID int) VersionCandidates {
	candidates := VersionCandidates{}
	for _, output := range db.ResourceVersions {
		if output.ResourceSpaceID == resourceID {
			candidates.Add(VersionCandidate{
				VersionID:  output.VersionID,
				CheckOrder: output.CheckOrder,
			})
		}
	}

	return candidates
}

func (db VersionsDB) LatestVersionOfResource(resourceID int) (VersionCandidate, bool) {
	var candidate VersionCandidate
	var found bool

	for _, v := range db.ResourceVersions {
		if v.ResourceSpaceID == resourceID && v.CheckOrder > candidate.CheckOrder {
			candidate = VersionCandidate{
				VersionID:  v.VersionID,
				CheckOrder: v.CheckOrder,
			}

			found = true
		}
	}

	return candidate, found
}

func (db VersionsDB) FindVersionOfResource(resourceID int, versionID int) (VersionCandidate, bool) {
	var candidate VersionCandidate
	var found bool

	for _, v := range db.ResourceVersions {
		if v.ResourceSpaceID == resourceID && v.VersionID == versionID {
			candidate = VersionCandidate{
				VersionID:  v.VersionID,
				CheckOrder: v.CheckOrder,
			}

			found = true
		}
	}

	return candidate, found
}

func (db VersionsDB) VersionsOfResourcePassedJobs(resourceID int, passed JobPermutationSet) VersionCandidates {
	candidates := VersionCandidates{}

	firstTick := true
	for jobID, _ := range passed {
		versions := VersionCandidates{}

		for _, output := range db.BuildOutputs {
			if output.ResourceSpaceID == resourceID && output.JobPermutationID == jobID {
				versions.Add(VersionCandidate{
					VersionID:        output.VersionID,
					CheckOrder:       output.CheckOrder,
					BuildID:          output.BuildID,
					JobPermutationID: output.JobPermutationID,
				})
			}
		}

		if firstTick {
			candidates = versions
			firstTick = false
		} else {
			candidates = candidates.IntersectByVersion(versions)
		}
	}

	return candidates
}
