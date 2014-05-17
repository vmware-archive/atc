package builds

import "encoding/json"

type Build struct {
	Guid string `json:"guid"`

	Privileged bool `json:"privileged"`

	Config Config `json:"config"`

	LogsURL  string `json:"logs_url"`
	Callback string `json:"callback"`

	Inputs []Input `json:"inputs"`

	Status string `json:"status"`
}

type Config struct {
	Image string `json:"image"`

	Env    [][2]string `json:"env"`
	Script string      `json:"script"`
}

type Input struct {
	Type string `json:"type"`

	ConfigPath      string `json:"configPath"`
	DestinationPath string `json:"destinationPath"`

	Version *json.RawMessage `json:"version"`
}
