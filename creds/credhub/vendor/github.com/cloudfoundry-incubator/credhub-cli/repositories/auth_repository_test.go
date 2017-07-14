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
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

var _ = Describe("AuthRepository", func() {
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

		Context("when there is a response body", func() {
			BeforeEach(func() {
				repository = NewAuthRepository(&httpClient, true)
			})

			It("sends a request to the server", func() {
				request, _ := http.NewRequest("POST", cfg.AuthURL, nil)

				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"access_token":"token","expires_in": 1234,"refresh_token":"refresh-token"}`))),
				}

				httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
					Expect(req).To(Equal(request))

					return &responseObj, nil
				}

				expectedToken := models.Token{
					AccessToken:  "token",
					ExpiresIn:    1234,
					RefreshToken: "refresh-token",
				}

				token, err := repository.SendRequest(request, "")

				Expect(err).ToNot(HaveOccurred())
				Expect(token).To(Equal(expectedToken))
			})

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

		})

		Context("when there is no response body", func() {

			BeforeEach(func() {
				repository = NewAuthRepository(&httpClient, false)
			})

			It("does not require a response body", func() {
				responseObj := http.Response{StatusCode: 200}
				httpClient.DoReturns(&responseObj, nil)
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				_, err := repository.SendRequest(request, "foo")
				Expect(err).ToNot(HaveOccurred())
			})

		})
	})

})
