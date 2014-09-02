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

	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`

	EventsCallback string `json:"events_callback"`
	StatusCallback string `json:"status_callback"`

	AbortURL  string `json:"abort_url"`
	HijackURL string `json:"hijack_url"`

	Status Status `json:"status"`
}

type Config struct {
	Image  string            `json:"image"   yaml:"image"`
	Params map[string]string `json:"params"  yaml:"params"`
	Run    RunConfig         `json:"run"     yaml:"run"`
	Paths  map[string]string `json:"paths"   yaml:"paths"`
}

type RunConfig struct {
	Path string   `json:"path" yaml:"path"`
	Args []string `json:"args" yaml:"args"`
}

type Input struct {
	Name string `json:"name"`

	Type string `json:"type"`

	// e.g. sha
	Version Version `json:"version,omitempty"`

	// e.g. git url, branch, private_key
	Source Source `json:"source"`

	// arbitrary config for input
	Params Params `json:"params,omitempty"`

	// e.g. commit_author, commit_date
	Metadata []MetadataField `json:"metadata,omitempty"`

	ConfigPath string `json:"config_path"`
}

type Version map[string]interface{}

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Output struct {
	Name string `json:"name"`

	Type string `json:"type"`

	// e.g. [success, failure]
	On []OutputCondition `json:"on"`

	// e.g. sha
	Version Version `json:"version"`

	// e.g. git url, branch, private_key
	Source Source `json:"source"`

	// arbitrary config for output
	Params Params `json:"params,omitempty"`

	// e.g. commit_author, commit_date, commit_sha
	Metadata []MetadataField `json:"metadata,omitempty"`
}

type OutputCondition string

const (
	OutputConditionSuccess OutputCondition = "success"
	OutputConditionFailure OutputCondition = "failure"
)

type Source map[string]interface{}

type Params map[string]interface{}
