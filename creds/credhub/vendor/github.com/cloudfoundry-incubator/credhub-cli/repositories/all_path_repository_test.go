package repositories_test

import (
	. "github.com/cloudfoundry-incubator/credhub-cli/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/client/clientfakes"
	"github.com/cloudfoundry-incubator/credhub-cli/config"

	credhub_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
)

var _ = Describe("FindRepository", func() {
	var (
		repository Repository
		httpClient clientfakes.FakeHttpClient
		cfg        config.Config
	)

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:  "http://example.com",
			AuthURL: "http://uaa.example.com",
		}
	})

	Describe("SendRequest", func() {
		BeforeEach(func() {
			repository = NewAllPathRepository(&httpClient)
		})

		It("sends a request to the server", func() {
			request, _ := http.NewRequest("GET", "http://example.com/data?paths=true", nil)

			// language=JSON
			expectedJson := `{
                "paths": [
                    {"path": "deploy123/"},
                    {"path": "deploy123/dan/"},
                    {"path": "deploy123/dan/consul/"},
                    {"path": "deploy12/"},
                    {"path": "consul/"},
                    {"path": "consul/deploy123/"}
                ]}`

			responseObj := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(expectedJson))),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			findResponseBody, err := repository.SendRequest(request, "")
			Expect(err).ToNot(HaveOccurred())

			Expect(findResponseBody.ToJson()).To(MatchJSON(expectedJson))
		})

		Describe("Errors", func() {
			It("returns a NewResponseError when the JSON response cannot be parsed", func() {
				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("adasdasdasasd"))),
				}
				httpClient.DoReturns(&responseObj, nil)
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				_, error := repository.SendRequest(request, "foo")
				Expect(error).To(MatchError(credhub_errors.NewResponseError()))
			})

			It("returns a NewNoMatchingCredentialsFoundError when there are no credentials returned", func() {
				emptyPaths := `{"paths":[]}`

				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(emptyPaths))),
				}
				httpClient.DoReturns(&responseObj, nil)
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				paths, err := repository.SendRequest(request, "foo")
				Expect(err).To(MatchError(credhub_errors.NewNoMatchingCredentialsFoundError()))
				Expect(paths.ToJson()).To(MatchJSON(emptyPaths))
			})
		})
	})

})
