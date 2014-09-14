package event

import "github.com/concourse/turbine/api/builds"

type Input struct {
	Input builds.Input `json:"input"`
}

func (Input) EventType() EventType { return EventTypeInput }
