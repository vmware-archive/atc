package algorithm

type VersionsDB struct {
	ResourceVersions  []ResourceVersion
	BuildOutputs      []BuildOutput
	BuildInputs       []BuildInput
	JobCombinationIDs map[string]int
	ResourceIDs       map[string]int
}

type ResourceVersion struct {
	VersionID  int
	ResourceID int
	CheckOrder int
}

type BuildOutput struct {
	ResourceVersion
	BuildID          int
	JobCombinationID int
}

type BuildInput struct {
	ResourceVersion
	BuildID          int
	JobCombinationID int
	InputName        string
}

func (db VersionsDB) IsVersionFirstOccurrence(versionID int, jobCombinationID int, inputName string) bool {
	for _, buildInput := range db.BuildInputs {
		if buildInput.VersionID == versionID &&
			buildInput.JobCombinationID == jobCombinationID &&
			buildInput.InputName == inputName {
			return false
		}
	}
	return true
}

func (db VersionsDB) AllVersionsOfResource(resourceID int) VersionCandidates {
	candidates := VersionCandidates{}
	for _, output := range db.ResourceVersions {
		if output.ResourceID == resourceID {
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
		if v.ResourceID == resourceID && v.CheckOrder > candidate.CheckOrder {
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
		if v.ResourceID == resourceID && v.VersionID == versionID {
			candidate = VersionCandidate{
				VersionID:  v.VersionID,
				CheckOrder: v.CheckOrder,
			}

			found = true
		}
	}

	return candidate, found
}

func (db VersionsDB) VersionsOfResourcePassedJobs(resourceID int, passed JobSet) VersionCandidates {
	candidates := VersionCandidates{}

	firstTick := true
	for jobCombinationID, _ := range passed {
		versions := VersionCandidates{}

		for _, output := range db.BuildOutputs {
			if output.ResourceID == resourceID && output.JobCombinationID == jobCombinationID {
				versions.Add(VersionCandidate{
					VersionID:        output.VersionID,
					CheckOrder:       output.CheckOrder,
					BuildID:          output.BuildID,
					JobCombinationID: output.JobCombinationID,
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
