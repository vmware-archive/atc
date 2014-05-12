package builds

type Build struct {
	Guid string `json:"guid"`

	ConfigPath string `json:"config"`

	Image  string      `json:"image"`
	Env    [][2]string `json:"env"`
	Script string      `json:"script"`

	Privileged bool `json:"privileged"`

	LogsURL  string `json:"logs_url"`
	Callback string `json:"callback"`

	Sources []BuildSource `json:"sources"`

	Status string `json:"status"`
}

type BuildSource struct {
	Type   string `json:"type"`
	URI    string `json:"uri"`
	Branch string `json:"branch"`
	Ref    string `json:"ref"`
	Path   string `json:"path"`
}
