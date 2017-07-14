package models

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CredentialResponse", func() {
	expectedJSON := `{"name": "stringSecret",
					"type":       "value",
					"value":            "my-value",
					"version_created_at": "2016-01-01T12:00:00Z"}`

	credential := CredentialResponse{
		make(map[string]interface{}),
	}

	Describe("ToYaml()", func() {
		It("renders string secrets", func() {
			err := json.Unmarshal([]byte(expectedJSON), &credential.ResponseBody)
			Expect(err).To(BeNil())

			Expect(credential.ToYaml()).To(MatchYAML("" +
				"type: value\n" +
				"name: stringSecret\n" +
				"value: my-value\n" +
				"version_created_at: 2016-01-01T12:00:00Z\n"))
		})
	})

	Describe("ToJson()", func() {
		It("renders string secrets", func() {
			err := json.Unmarshal([]byte(expectedJSON), &credential.ResponseBody)
			Expect(err).To(BeNil())
			Expect(credential.ToJson()).To(MatchJSON(expectedJSON))
		})
	})
})
