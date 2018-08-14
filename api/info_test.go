package api_test

import (
	"io/ioutil"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/aws/aws-sdk-go/aws"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/concourse/atc/api/accessor/accessorfakes"
	"github.com/concourse/atc/creds/credhub"
	"github.com/concourse/atc/creds/ssm"
	"github.com/concourse/atc/creds/vault"
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

type mockedSsmClient struct {
	ssmiface.SSMAPI
}

type MockSsmService struct {
	ssmiface.SSMAPI

	stubDescribeInstanceInformation func(input *awsssm.DescribeInstanceInformationInput) (awsssm.DescribeInstanceInformationOutput, error)
}

var _ = Describe("Pipelines API", func() {
	Describe("GET /api/v1/info", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			response, err = client.Get(server.URL + "/api/v1/info")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns Content-Type 'application/json'", func() {
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json"))
		})

		It("contains the version", func() {
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).To(MatchJSON(`{
				"version": "1.2.3",
				"worker_version": "4.5.6"
			}`))
		})
	})

	Describe("GET /api/v1/info/creds", func() {
		var (
			response   *http.Response
			fakeaccess *accessorfakes.FakeAccess
			credServer *ghttp.Server
		)

		BeforeEach(func() {
			fakeaccess = new(accessorfakes.FakeAccess)
		})

		JustBeforeEach(func() {
			fakeAccessor.CreateReturns(fakeaccess)
		})

		JustBeforeEach(func() {
			var err error

			req, err := http.NewRequest("GET", server.URL+"/api/v1/info/creds", nil)
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			response, err = client.Do(req)
			Expect(err).NotTo(HaveOccurred())
		})

		FContext("SSM", func() {
			var mockService MockSsmService
			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(true)
				fakeaccess.IsAdminReturns(true)

				ssmManager := &ssm.SsmManager{
					AwsAccessKeyID:         "",
					AwsSecretAccessKey:     "",
					AwsSessionToken:        "",
					AwsRegion:              "blah",
					PipelineSecretTemplate: "",
					TeamSecretTemplate:     "",
					Client:                 &mockedSsmClient{},
				}

				credsManagers["ssm"] = ssmManager

				t1, err := template.New("test").Parse(DefaultPipelineSecretTemplate)
				Expect(t1).NotTo(BeNil())
				Expect(err).To(BeNil())
				t2, err := template.New("test").Parse(DefaultTeamSecretTemplate)
				Expect(t2).NotTo(BeNil())
				Expect(err).To(BeNil())
				ssmAccess = NewSsm(lager.NewLogger("ssm_test"), &mockService, "alpha", "bogus", []*template.Template{t1, t2})
				mockService.stubDescribeInstanceInformation = func(input *awsssm.DescribeInstanceInformationInput) (awsssm.DescribeInstanceInformationOutput, error) {
					return awsssm.DescribeInstanceInformationOutput{
						PingStatus: aws.String("Online"),
					}
				}
			})

			It("returns configured ssm manager", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
          "ssm": {
						"region": "blah",
						"health": {
							"InstanceInformationList": [
								{
									"PingStatus": "Online"
								}
							]
						}
          }
        }`))
			})
		})

		Context("vault", func() {
			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(true)
				fakeaccess.IsAdminReturns(true)

				authConfig := vault.AuthConfig{
					Backend:       "backend-server",
					BackendMaxTTL: 20,
					RetryMax:      5,
					RetryInitial:  2,
				}

				tls := vault.TLS{
					CACert:     "",
					ServerName: "server-name",
				}

				credServer = ghttp.NewServer()
				vaultManager := &vault.VaultManager{
					URL:        credServer.URL(),
					PathPrefix: "testpath",
					Cache:      false,
					MaxLease:   60,
					TLS:        tls,
					Auth:       authConfig,
				}

				credsManagers["vault"] = vaultManager

				credServer.RouteToHandler("GET", "/v1/sys/health", ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					&vaultapi.HealthResponse{
						Initialized:                true,
						Sealed:                     false,
						Standby:                    false,
						ReplicationPerformanceMode: "foo",
						ReplicationDRMode:          "blah",
						ServerTimeUTC:              0,
						Version:                    "1.0.0",
					},
				))
			})

			It("returns Content-Type 'application/json'", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json"))
			})

			It("returns configured creds manager", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
          "vault": {
            "url": "` + credServer.URL() + `",
            "path_prefix": "testpath",
						"cache": false,
						"max_lease": 60,
            "ca_cert": "",
            "server_name": "server-name",
						"auth_backend": "backend-server",
						"auth_max_ttl": 20,
						"auth_retry_max": 5,
						"auth_retry_initial": 2,
						"health": {
							"initialized": true,
							"sealed": false,
							"standby": false,
							"replication_performance_mode": "foo",
							"replication_dr_mode": "blah",
							"server_time_utc": 0,
							"version": "1.0.0"
						}
          }
        }`))
			})
		})

		Context("credhub", func() {
			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(true)
				fakeaccess.IsAdminReturns(true)

				tls := credhub.TLS{
					CACerts: []string{},
				}
				uaa := credhub.UAA{
					ClientId:     "client-id",
					ClientSecret: "client-secret",
				}
				credServer = ghttp.NewServer()
				credhubManager := &credhub.CredHubManager{
					URL:        credServer.URL(),
					PathPrefix: "some-prefix",
					TLS:        tls,
					UAA:        uaa,
				}

				credsManagers["credhub"] = credhubManager
				credServer.RouteToHandler("GET", "/health", ghttp.RespondWithJSONEncoded(
					http.StatusOK, map[string]string{
						"status": "UP",
					},
				))
				credServer.RouteToHandler("GET", "/info", ghttp.RespondWithJSONEncoded(
					http.StatusOK, map[string]interface{}{
						"auth-server": map[string]string{"url": "http://1.2.3.4:8080"},
						"app":         map[string]string{"app": "CredHub"},
					},
				))

			})

			It("returns Content-Type 'application/json'", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json"))
			})

			It("returns configured creds manager", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
          "credhub": {
            "url": "` + credServer.URL() + `",
            "path_prefix": "some-prefix",
            "ca_certs": [],
						"uaa_client_id": "client-id",
						"health": {"status": "UP"}
          }
        }`))
			})
		})

	})
})
