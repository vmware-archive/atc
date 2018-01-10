package algorithm

type ExistingBuildResolver struct {
	BuildInputs      []BuildInput
	JobCombinationID int
	ResourceID       int
}

func (r *ExistingBuildResolver) Exists() bool {
	for _, buildInput := range r.BuildInputs {
		if buildInput.JobCombinationID == r.JobCombinationID && buildInput.ResourceID == r.ResourceID {
			return true
		}
	}

	return false
}

func (r *ExistingBuildResolver) ExistsForVersion(versionID int) bool {
	for _, buildInput := range r.BuildInputs {
		if buildInput.JobCombinationID == r.JobCombinationID && buildInput.ResourceID == r.ResourceID {
			if buildInput.VersionID == versionID {
				return true
			}
		}
	}

	return false
}
