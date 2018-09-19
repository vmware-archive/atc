package present

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

func ResourceConfigVersion(o db.BuildOutput) atc.VersionedResource {
	return atc.VersionedResource{
		Resource: o.Name,
		Version:  atc.Version(o.Version),
	}
}
