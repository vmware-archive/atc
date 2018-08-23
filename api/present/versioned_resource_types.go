package present

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

func VersionedResourceTypes(showCheckError bool, savedResourceTypes db.ResourceTypes) atc.VersionedResourceTypes {
	versionedResourceTypes := savedResourceTypes.Deserialize()

	for i, resourceType := range savedResourceTypes {
		if resourceType.CheckError() != nil && showCheckError {
			versionedResourceTypes[i].CheckSetupError = resourceType.CheckError().Error()
		} else {
			versionedResourceTypes[i].CheckSetupError = ""
		}

		if resourceType.ResourceConfigCheckError() != nil && showCheckError {
			versionedResourceTypes[i].CheckError = resourceType.ResourceConfigCheckError().Error()
		} else {
			versionedResourceTypes[i].CheckError = ""
		}
	}

	return versionedResourceTypes
}
