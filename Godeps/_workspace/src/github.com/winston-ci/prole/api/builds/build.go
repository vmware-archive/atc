package builds

import "fmt"

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

	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`

	Status Status `json:"status"`
}

type Config struct {
	Image string `json:"image"`

	Env    [][2]string `json:"env"`
	Script string      `json:"script"`
}

type Input struct {
	Name string `json:"name"`

	Type   string `json:"type"`
	Source Source `json:"source,omitempty"`

	ConfigPath      string `json:"config_path"`
	DestinationPath string `json:"destination_path"`

	Metadata []MetadataField `json:"metadata"`
}

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Source []byte

func (source Source) MarshalJSON() ([]byte, error) {
	return []byte(source), nil
}

func (source *Source) UnmarshalJSON(data []byte) error {
	*source = append((*source)[0:0], data...)
	return nil
}

func (source Source) String() string {
	return string(source)
}

func (source Source) GoString() string {
	return fmt.Sprintf("builds.Source(%q)", source)
}

type Output struct {
	Name string `json:"name"`

	Type   string `json:"type"`
	Params Params `json:"params,omitempty"`

	Source Source `json:"source,omitempty"`

	SourcePath string `json:"sourcePath"`
}

type Params []byte

func (params Params) MarshalJSON() ([]byte, error) {
	return []byte(params), nil
}

func (params *Params) UnmarshalJSON(data []byte) error {
	*params = append((*params)[0:0], data...)
	return nil
}

func (params Params) String() string {
	return string(params)
}

func (params Params) GoString() string {
	return fmt.Sprintf("builds.Params(%q)", params)
}
