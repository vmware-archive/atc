package event

type Start struct {
	Time int64 `json:"time"`
}

func (Start) EventType() EventType { return EventTypeStart }
