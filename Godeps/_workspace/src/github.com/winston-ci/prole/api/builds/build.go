package builds

type Status string

const (
	StatusStarted   Status = "started"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
	StatusErrored   Status = "errored"
)

type Build struct {
	Guid string `json:"guid"`

	Privileged bool `json:"privileged"`

	Config Config `json:"config"`

	LogsURL  string `json:"logs_url"`
	Callback string `json:"callback"`

	Inputs []Input `json:"inputs"`

	Status Status `json:"status"`
}

type Config struct {
	Image string `json:"image"`

	Env    [][2]string `json:"env"`
	Script string      `json:"script"`
}

type Source []byte

func (source Source) MarshalJSON() ([]byte, error) {
	return []byte(source), nil
}

func (source *Source) UnmarshalJSON(data []byte) error {
	*source = append((*source)[0:0], data...)
	return nil
}

type Input struct {
	Type string `json:"type"`

	Source Source `json:"source,omitempty"`

	ConfigPath      string `json:"configPath"`
	DestinationPath string `json:"destinationPath"`
}
