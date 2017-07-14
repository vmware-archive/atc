package models

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type CredentialResponse struct {
	ResponseBody map[string]interface{}
}

func (response CredentialResponse) ToYaml() string {
	s, _ := yaml.Marshal(response.ResponseBody)
	return string(s)
}

func (response CredentialResponse) ToJson() string {
	s, _ := json.MarshalIndent(response.ResponseBody, "", "\t")
	return string(s)
}
