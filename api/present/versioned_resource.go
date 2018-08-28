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

func SavedVersionedResource(svr db.SavedVersionedResource) atc.VersionedResource {
	var metadata []atc.MetadataField

	for _, v := range svr.Metadata {
		metadata = append(metadata, atc.MetadataField(v))
	}

	return atc.VersionedResource{
		ID:       svr.ID,
		Resource: svr.Resource,
		Enabled:  svr.Enabled,
		Type:     svr.Type,
		Version:  atc.Version(svr.Version),
		Metadata: metadata,
	}
}
