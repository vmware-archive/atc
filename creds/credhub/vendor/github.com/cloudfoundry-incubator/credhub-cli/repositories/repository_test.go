package repositories_test

import (
	. "github.com/cloudfoundry-incubator/credhub-cli/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io/ioutil"
	"net/http"

	"errors"

	"github.com/cloudfoundry-incubator/credhub-cli/client/clientfakes"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	credhub_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
)

var _ = Describe("Repository", func() {
	var (
		httpClient clientfakes.FakeHttpClient
		cfg        config.Config
	)

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:  "http://example.com",
			AuthURL: "http://uaa.example.com",
		}
	})

	Describe("DoSendRequest", func() {
		It("sends a request to the server", func() {
			request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

			responseObj := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"type":"value","value":"my-value"}`))),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			response, err := DoSendRequest(&httpClient, request)
			Expect(err).ToNot(HaveOccurred())
			Expect(response).To(Equal(&responseObj))
		})

		Describe("Errors", func() {
			It("returns NetworkError when there is a network error", func() {
				serverError := errors.New("hello")
				httpClient.DoReturns(nil, serverError)

				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)
				_, error := DoSendRequest(&httpClient, request)
				Expect(error).To(MatchError(credhub_errors.NewNetworkError(serverError)))
			})

			It("returns a error when response is 400", func() {
				responseObj := http.Response{
					StatusCode: 400,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "My error"}`))),
				}

				httpClient.DoReturns(&responseObj, nil)

				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)
				_, error := DoSendRequest(&httpClient, request)

				Expect(error.Error()).To(Equal("My error"))
			})

			It("returns UnauthorizedError when the CM server returns Unauthorized", func() {
				responseObj := http.Response{
					StatusCode: 401,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "invalid_token","error_description":"My long error text"}`))),
				}
				httpClient.DoReturns(&responseObj, nil)
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				_, error := DoSendRequest(&httpClient, request)
				Expect(error).To(MatchError(errors.New("My long error text")))
			})

			It("returns ForbiddenError when the CM server returns Forbidden", func() {
				responseObj := http.Response{
					StatusCode: 403,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "insufficient_scope","error_description":"Insufficient scope for this resource"}`))),
				}
				httpClient.DoReturns(&responseObj, nil)
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				_, error := DoSendRequest(&httpClient, request)
				Expect(error).To(MatchError(credhub_errors.NewForbiddenError()))
			})

			It("returns an error when response is 500", func() {
				responseObj := http.Response{
					StatusCode: 500,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "My error"}`))),
				}

				httpClient.DoReturns(&responseObj, nil)

				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)
				_, error := DoSendRequest(&httpClient, request)

				Expect(error.Error()).To(Equal("My error"))
			})

			It("returns generic 500 error when response is 500 and there is no body response from the server", func() {
				responseObj := http.Response{
					StatusCode: 500,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(``))),
				}

				httpClient.DoReturns(&responseObj, nil)

				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)
				_, error := DoSendRequest(&httpClient, request)

				Expect(error.Error()).NotTo(Equal("EOF"))
				Expect(error.Error()).To(Equal("The targeted API was unable to perform the request. Please validate and retry your request."))
			})

			It("returns AccessTokenExpiredError when server indicates the token has expired", func() {
				responseObj := http.Response{
					StatusCode: 401,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "access_token_expired","error_description":"More long error text"}`))),
				}
				httpClient.DoReturns(&responseObj, nil)
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				_, error := DoSendRequest(&httpClient, request)
				Expect(error).To(MatchError(credhub_errors.NewAccessTokenExpiredError()))
			})
		})
	})
})
