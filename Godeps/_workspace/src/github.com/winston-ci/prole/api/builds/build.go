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

	AbortURL string `json:"abort_url"`

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

	Env    []map[string]string `json:"env"`
	Script string              `json:"script"`
}

type Input struct {
	Name string `json:"name"`

	Type string `json:"type"`

	// e.g. sha
	Version Version `json:"version,omitempty"`

	// e.g. git url, branch, private_key
	Source Source `json:"source"`

	// e.g. commit_author, commit_date
	Metadata []MetadataField `json:"metadata,omitempty"`

	ConfigPath      string `json:"config_path"`
	DestinationPath string `json:"destination_path"`
}

type Version map[string]interface{}

type Source map[string]interface{}

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Output struct {
	Name string `json:"name"`

	Type string `json:"type"`

	// e.g. sha
	Version Version `json:"version"`

	// e.g. git url, branch
	Params Params `json:"params,omitempty"`

	// e.g. commit_author, commit_date, commit_sha
	Metadata []MetadataField `json:"metadata,omitempty"`

	SourcePath string `json:"sourcePath"`
}

type Params map[string]interface{}
