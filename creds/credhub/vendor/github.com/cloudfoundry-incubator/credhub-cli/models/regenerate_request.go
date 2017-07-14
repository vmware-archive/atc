package models

type RegenerateRequest struct {
	Name       string `json:"name"`
	Regenerate bool   `json:"regenerate"`
}
