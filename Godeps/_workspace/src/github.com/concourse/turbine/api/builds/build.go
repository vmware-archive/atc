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

	Privileged bool   `json:"privileged"`
	Config     Config `json:"config"`

	AbortURL string `json:"abort_url"`
	LogsURL  string `json:"logs_url"`
	Callback string `json:"callback"`

	Status Status `json:"status"`
}

type Config struct {
	Image   string            `json:"image"   yaml:"image"`
	Params  map[string]string `json:"params"  yaml:"params"`
	Run     RunConfig         `json:"run"     yaml:"run"`
	Inputs  []Input           `json:"inputs"  yaml:"inputs"`
	Outputs []Output          `json:"outputs" yaml:"outputs"`
}

type RunConfig struct {
	Path string   `json:"path" yaml:"path"`
	Args []string `json:"args" yaml:"args"`
}

type Input struct {
	Name string `json:"name" yaml:"name"`

	Type string `json:"type" yaml:"type"`

	// e.g. sha
	Version Version `json:"version,omitempty" yaml:"version"`

	// e.g. git url, branch, private_key
	Source Source `json:"source" yaml:"source"`

	// e.g. commit_author, commit_date
	Metadata []MetadataField `json:"metadata,omitempty" yaml"metadata"`

	ConfigPath string `json:"config_path" yaml:"config_path"`

	DestinationPath string `json:"destination_path" yaml:"destination_path"`
}

type Version map[string]interface{}

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Output struct {
	Name string `json:"name"`

	Type string `json:"type"`

	// e.g. sha
	Version Version `json:"version"`

	// e.g. git url, branch, private_key
	Source Source `json:"source"`

	// arbitrary config for output
	Params Params `json:"params,omitempty"`

	// e.g. commit_author, commit_date, commit_sha
	Metadata []MetadataField `json:"metadata,omitempty"`

	SourcePath string `json:"sourcePath"`
}

type Source map[string]interface{}

type Params map[string]interface{}
