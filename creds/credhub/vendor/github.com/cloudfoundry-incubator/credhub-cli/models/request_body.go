package models

type RequestBody struct {
	CredentialType   string                `json:"type" binding:"required"`
	Name             string                `json:"name,omitempty"`
	Value            interface{}           `json:"value,omitempty"`
	Overwrite        *bool                 `json:"overwrite,omitempty"`
	Parameters       *GenerationParameters `json:"parameters,omitempty"`
	VersionCreatedAt string                `json:"version_created_at,omitempty"`
}
