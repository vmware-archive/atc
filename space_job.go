package atc

type SpaceJob struct {
	ID int `json:"id"`

	Name                 string `json:"name"`
	PipelineName         string `json:"pipeline_name"`
	TeamName             string `json:"team_name"`
	Paused               bool   `json:"paused,omitempty"`
	FirstLoggedBuildID   int    `json:"first_logged_build_id,omitempty"`
	DisableManualTrigger bool   `json:"disable_manual_trigger,omitempty"`

	Inputs  []JobInput  `json:"inputs"`
	Outputs []JobOutput `json:"outputs"`

	Groups []string `json:"groups"`

	Combinations []SpaceJobCombination `json:"combinations"`
}

type SpaceJobCombination struct {
	ID          int               `json:"id"`
	Combination map[string]string `json:"combination"`

	NextBuild       *Build `json:"next_build"`
	FinishedBuild   *Build `json:"finished_build"`
	TransitionBuild *Build `json:"transition_build,omitempty"`
}
