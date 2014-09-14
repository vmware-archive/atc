package event

type Error struct {
	Message string `json:"message"`
}

func (Error) EventType() EventType { return EventTypeError }
