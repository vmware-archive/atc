package event

import "github.com/concourse/turbine/api/builds"

type Initialize struct {
	BuildConfig builds.Config `json:"config"`
}

func (Initialize) EventType() EventType { return EventTypeInitialize }
