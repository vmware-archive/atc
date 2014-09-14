package event

import "github.com/concourse/turbine/api/builds"

type Status struct {
	Status builds.Status `json:"status"`
}

func (Status) EventType() EventType { return EventTypeStatus }
