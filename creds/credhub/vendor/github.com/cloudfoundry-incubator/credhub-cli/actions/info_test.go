package actions_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/client/clientfakes"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info", func() {
	var (
		subject    actions.ServerInfo
		httpClient clientfakes.FakeHttpClient
		cfg        config.Config
	)

	BeforeEach(func() {
		cfg = config.Config{ApiURL: "example.com"}
		subject = actions.NewInfo(&httpClient, cfg)
	})

	Describe("ServerInfo", func() {
		It("returns the version of the cli and CM server, as well as auth server URL", func() {
			request := client.NewInfoRequest(cfg)

			responseObj := http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"app":{"version":"my-version","name":"CredHub"},
					"auth-server":{"url":"https://example.com"}
					}`)),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			serverInfo, _ := subject.GetServerInfo()
			Expect(serverInfo.App.Version).To(Equal("my-version"))
			Expect(serverInfo.AuthServer.Url).To(Equal("https://example.com"))
		})

		It("returns error if server returned a non 200 status code", func() {
			responseObj := http.Response{StatusCode: 400}

			httpClient.DoReturns(&responseObj, nil)

			_, err := subject.GetServerInfo()
			Expect(err).NotTo(BeNil())
		})

		It("returns error if server has a network error", func() {
			responseObj := http.Response{
				StatusCode: 200,
			}

			httpClient.DoReturns(&responseObj, errors.New("dogs are gone"))

			_, err := subject.GetServerInfo()
			Expect(err).NotTo(BeNil())
		})

		It("returns error if server returns bad json", func() {
			responseObj := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`sdafasdfasdf`)),
			}

			httpClient.DoReturns(&responseObj, nil)

			_, err := subject.GetServerInfo()
			Expect(err).NotTo(BeNil())
		})
	})
})
