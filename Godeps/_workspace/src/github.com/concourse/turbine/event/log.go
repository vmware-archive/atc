package event

type Log struct {
	Origin  Origin `json:"origin"`
	Payload string `json:"payload"`
}

func (Log) EventType() EventType { return EventTypeLog }

type Origin struct {
	Type OriginType `json:"type"`
	Name string     `json:"name"`
}

type OriginType string

const (
	OriginTypeInvalid OriginType = ""
	OriginTypeInput   OriginType = "input"
	OriginTypeOutput  OriginType = "output"
	OriginTypeRun     OriginType = "run"
)
