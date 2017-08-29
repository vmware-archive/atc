package client_test

import (
	"net/http"

	. "github.com/cloudfoundry-incubator/credhub-cli/client"

	"bytes"

	"fmt"

	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Credhub API", func() {
	var cfg config.Config

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:      "http://example.com",
			AuthURL:     "http://example.com/uaa",
			AccessToken: "access-token",
		}
	})

	Describe("NewInfoRequest", func() {
		It("Returns a request for the info endpoint", func() {
			expectedRequest, _ := http.NewRequest("GET", "http://example.com/info", nil)

			request := NewInfoRequest(cfg)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewSetCredentialRequest with value", func() {
		It("Returns a request for the put-value endpoint", func() {
			request := NewSetCredentialRequest(cfg, "value", "my-name", "my-value", true)

			Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
			Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer access-token"}))
			Expect(request.Method).To(Equal("PUT"))

			byteBuff := new(bytes.Buffer)
			byteBuff.ReadFrom(request.Body)

			Expect(byteBuff.String()).To(MatchJSON(`{"type":"value","name":"my-name","value":"my-value","overwrite":true}`))
		})

		It("Returns a request that will not overwrite", func() {
			request := NewSetCredentialRequest(cfg, "value", "my-name", "my-value", false)

			Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
			Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer access-token"}))
			Expect(request.Method).To(Equal("PUT"))

			byteBuff := new(bytes.Buffer)
			byteBuff.ReadFrom(request.Body)

			Expect(byteBuff.String()).To(MatchJSON(`{"type":"value","name":"my-name","value":"my-value","overwrite":false}`))
		})
	})

	Describe("NewSetCredentialRequest with password", func() {
		It("Returns a request for the put-password endpoint", func() {
			request := NewSetCredentialRequest(cfg, "password", "my-name", "my-password", true)
			Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
			Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer access-token"}))
			Expect(request.Method).To(Equal("PUT"))

			byteBuff := new(bytes.Buffer)
			byteBuff.ReadFrom(request.Body)

			Expect(byteBuff.String()).To(MatchJSON(`{"type":"password","name":"my-name","value":"my-password","overwrite":true}`))
		})
	})

	Describe("NewPutCertificateRequest", func() {
		It("Returns a request for the put-certificate endpoint", func() {
			json := fmt.Sprintf(`{"type":"certificate","name":"my-name","value":{"ca":"%s","certificate":"%s","private_key":"%s"},"overwrite":true}`,
				"my-ca", "my-cert", "my-priv")

			request := NewSetCertificateRequest(cfg, "my-name", "my-ca", "", "my-cert", "my-priv", true)
			Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
			Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer access-token"}))
			Expect(request.Method).To(Equal("PUT"))

			byteBuff := new(bytes.Buffer)
			byteBuff.ReadFrom(request.Body)

			Expect(byteBuff.String()).To(MatchJSON(json))
		})
	})

	Describe("NewPutRsaSshRequest", func() {
		Describe("of type SSH", func() {
			It("Returns a request for the put-rsa-ssh endpoint", func() {
				json := fmt.Sprintf(`{"type":"%s","name":"my-name","value":{"public_key":"%s","private_key":"%s"},"overwrite":true}`,
					"ssh", "my-pub", "my-priv")

				request := NewSetRsaSshRequest(cfg, "my-name", "ssh", "my-pub", "my-priv", true)

				Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer access-token"}))
				Expect(request.Method).To(Equal("PUT"))

				byteBuff := new(bytes.Buffer)
				byteBuff.ReadFrom(request.Body)

				Expect(byteBuff.String()).To(MatchJSON(json))
			})
		})

		Describe("of type RSA", func() {
			It("Returns a request for the put-rsa-ssh endpoint", func() {
				json := fmt.Sprintf(`{"type":"%s","name":"my-name","value":{"public_key":"%s","private_key":"%s"},"overwrite":true}`,
					"rsa", "my-pub", "my-priv")

				request := NewSetRsaSshRequest(cfg, "my-name", "rsa", "my-pub", "my-priv", true)

				Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer access-token"}))
				Expect(request.Method).To(Equal("PUT"))

				byteBuff := new(bytes.Buffer)
				byteBuff.ReadFrom(request.Body)

				Expect(byteBuff.String()).To(MatchJSON(json))
			})
		})
	})

	Describe("NewPutUserRequest", func() {
		It("Returns a request for the put-user endpoint", func() {
			json := fmt.Sprintf(`{"type":"user","name":"my-name","value":{"username":"%s","password":"%s"},"overwrite":true}`,
				"my-user", "my-password")

			request := NewSetUserRequest(cfg, "my-name", "my-user", "my-password", true)
			Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
			Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer access-token"}))
			Expect(request.Method).To(Equal("PUT"))

			byteBuff := new(bytes.Buffer)
			byteBuff.ReadFrom(request.Body)

			Expect(byteBuff.String()).To(MatchJSON(json))
		})
	})

	Describe("NewGenerateCredentialRequest", func() {
		It("returns a request with only overwrite", func() {
			requestBody := `{"name":"my-name","type":"my-type","overwrite":true,"parameters":{}}`

			params := models.GenerationParameters{}
			request := NewGenerateCredentialRequest(cfg, "my-name", params, nil, "my-type", true)
			Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
			Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer access-token"}))
			Expect(request.Method).To(Equal("POST"))

			byteBuff := new(bytes.Buffer)
			byteBuff.ReadFrom(request.Body)

			Expect(byteBuff.String()).To(MatchJSON(requestBody))
		})

		It("returns a request with parameters", func() {
			parameters := models.GenerationParameters{
				IncludeSpecial: true,
				ExcludeNumber:  true,
				ExcludeUpper:   true,
				ExcludeLower:   true,
				Length:         42,
			}
			value := models.ProvidedValue{Username: "my-username"}
			expectedRequestBody := `{
					"name":"my-name",
					"type":"password",
					"overwrite":false,
					"parameters": {
						"include_special": true,
						"exclude_number": true,
						"exclude_upper": true,
						"exclude_lower": true,
						"length": 42
					},
					"value": {
						"username": "my-username"
					}
				}`

			request := NewGenerateCredentialRequest(cfg, "my-name", parameters, &value, "password", false)

			bodyBuffer := new(bytes.Buffer)
			bodyBuffer.ReadFrom(request.Body)
			Expect(bodyBuffer).To(MatchJSON(expectedRequestBody))
			Expect(request.Method).To(Equal("POST"))
			Expect(request.URL.String()).To(Equal("http://example.com/api/v1/data"))
			Expect(request.Header.Get("Content-Type")).To(Equal("application/json"))
		})
	})

	Describe("NewRegenerateCredentialRequest", func() {
		It("returns a request with only regenerate", func() {
			requestBody := `{"name":"my-name","regenerate":true}`

			request := NewRegenerateCredentialRequest(cfg, "my-name")
			Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
			Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer access-token"}))
			Expect(request.Method).To(Equal("POST"))

			byteBuff := new(bytes.Buffer)
			byteBuff.ReadFrom(request.Body)

			Expect(byteBuff.String()).To(MatchJSON(requestBody))
		})
	})

	Describe("NewGetCredentialByNameRequest", func() {
		It("Returns a request for getting a secret by name", func() {
			expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?name=my-name&current=true", nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewGetCredentialByNameRequest(cfg, "my-name")

			Expect(request).To(Equal(expectedRequest))
		})

		It("handles special characters in the query string", func() {
			rawName := "!wayt1cket/t0/cr@zy[town]?=AC/DC"
			escapedName := url.QueryEscape(rawName)

			Expect(escapedName).To(Equal("%21wayt1cket%2Ft0%2Fcr%40zy%5Btown%5D%3F%3DAC%2FDC"))

			expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?name="+escapedName+"&current=true", nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewGetCredentialByNameRequest(cfg, rawName)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewGetCredentialByIdRequest", func() {
		It("Returns a request for getting a secret by ID", func() {
			expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data/fake-test-id-123", nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewGetCredentialByIdRequest(cfg, "fake-test-id-123")

			Expect(request).To(Equal(expectedRequest))
		})

		It("handles special characters in the query string", func() {
			rawId := "!wayt1cket/t0/cr@zy[town]?=AC/DC"
			escapedId := url.QueryEscape(rawId)

			Expect(escapedId).To(Equal("%21wayt1cket%2Ft0%2Fcr%40zy%5Btown%5D%3F%3DAC%2FDC"))

			expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data/"+escapedId, nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewGetCredentialByIdRequest(cfg, rawId)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewFindCredentialsBySubstringRequest", func() {
		It("Returns a request for getting a credential", func() {
			expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?name-like=my-name", nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewFindCredentialsBySubstringRequest(cfg, "my-name")

			Expect(request).To(Equal(expectedRequest))
		})

		It("handles special characters in the query string", func() {
			rawName := "!wayt1cket/t0/cr@zy[town]?=AC/DC"
			escapedName := url.QueryEscape(rawName)

			expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?name-like="+escapedName, nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewFindCredentialsBySubstringRequest(cfg, rawName)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewFindAllCredentialPathsRequest", func() {
		It("Returns a request for getting all credential paths", func() {
			expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?paths=true", nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewFindAllCredentialPathsRequest(cfg)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewFindCredentialsByPathRequest", func() {
		It("Returns a request for getting a credential", func() {
			expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?path=my-path", nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewFindCredentialsByPathRequest(cfg, "my-path")

			Expect(request).To(Equal(expectedRequest))
		})

		It("handles special characters in the query string", func() {
			rawName := "!wayt1cket/t0/cr@zy[town]?=AC/DC"
			escapedName := url.QueryEscape(rawName)

			expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?path="+escapedName, nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewFindCredentialsByPathRequest(cfg, rawName)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewDeleteCredentialRequest", func() {
		It("Returns a request for deleting", func() {
			expectedRequest, _ := http.NewRequest("DELETE", "http://example.com/api/v1/data?name=my-name", nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewDeleteCredentialRequest(cfg, "my-name")

			Expect(request).To(Equal(expectedRequest))
		})

		It("handles special characters", func() {
			rawName := "?testParam=foo&gunk=x/bar/piv0t@l"
			escapedName := url.QueryEscape(rawName)

			expectedRequest, _ := http.NewRequest("DELETE", "http://example.com/api/v1/data?name="+escapedName, nil)
			expectedRequest.Header.Set("Authorization", "Bearer access-token")

			request := NewDeleteCredentialRequest(cfg, rawName)

			Expect(request).To(Equal(expectedRequest))
		})
	})
})
