package algorithm

type ExistingBuildResolver struct {
	BuildInputs      []BuildInput
	JobPermutationID int
	ResourceSpaceID  int
}

func (r *ExistingBuildResolver) Exists() bool {
	for _, buildInput := range r.BuildInputs {
		if buildInput.JobPermutationID == r.JobPermutationID && buildInput.ResourceSpaceID == r.ResourceSpaceID {
			return true
		}
	}

	return false
}

func (r *ExistingBuildResolver) ExistsForVersion(versionID int) bool {
	for _, buildInput := range r.BuildInputs {
		if buildInput.JobPermutationID == r.JobPermutationID && buildInput.ResourceSpaceID == r.ResourceSpaceID {
			if buildInput.VersionID == versionID {
				return true
			}
		}
	}

	return false
}
