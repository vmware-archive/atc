package event

type Finish struct {
	Time       int64 `json:"time"`
	ExitStatus int   `json:"exit_status"`
}

func (Finish) EventType() EventType { return EventTypeFinish }
