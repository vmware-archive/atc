package acceptance_test

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/lib/pq"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/concourse/atc/db"
)

var _ = Describe("Forcing HTTPS", func() {
	var atcProcess ifrit.Process
	var dbListener *pq.Listener
	var atcPort uint16

	BeforeEach(func() {
		atcBin, err := gexec.Build("github.com/concourse/atc/cmd/atc")
		Ω(err).ShouldNot(HaveOccurred())

		dbLogger := lagertest.NewTestLogger("test")
		postgresRunner.CreateTestDB()
		dbConn = postgresRunner.Open()
		dbListener = pq.NewListener(postgresRunner.DataSourceName(), time.Second, time.Minute, nil)
		bus := db.NewNotificationsBus(dbListener)
		sqlDB = db.NewSQL(dbLogger, dbConn, bus)

		var atcCommand *exec.Cmd
		atcCommand, atcPort = createATCCommandWithFlags(
			atcBin,
			1,
			map[string]string{
				"-forceHTTPS": "true",
			})
		atcProcess = startATC(atcCommand)
	})

	AfterEach(func() {
		ginkgomon.Interrupt(atcProcess)

		Ω(dbConn.Close()).Should(Succeed())
		Ω(dbListener.Close()).Should(Succeed())

		postgresRunner.DropTestDB()
	})

	Describe("making an API request", func() {
		var apiHomeURL string

		BeforeEach(func() {
			apiHomeURL = fmt.Sprintf("http://127.0.0.1:%d/api/v1/builds", atcPort)
		})

		Context("without forwarding header", func() {
			It("is rejected", func() {
				resp, err := http.Get(apiHomeURL)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
			})
		})

		Context("with forwarding header", func() {
			var client *http.Client
			var req *http.Request

			BeforeEach(func() {
				client = &http.Client{}
				req, _ = http.NewRequest("GET", apiHomeURL, nil)
				req.Header.Set("X-Forwarded-Proto", "https")
			})

			It("is accepted", func() {
				resp, err := client.Do(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("making a web request", func() {
		var webBuildsURL string
		var httpWebHomeURL string
		var host string
		const path = "/builds"

		BeforeEach(func() {
			host = fmt.Sprintf("127.0.0.1:%d", atcPort)
			webBuildsURL = fmt.Sprintf("%s%s", host, path)
			httpWebHomeURL = fmt.Sprintf("http://%s", webBuildsURL)
		})

		Context("without forwarding header", func() {
			var defaultTransport http.RoundTripper
			var req *http.Request

			BeforeEach(func() {
				defaultTransport = &http.Transport{}

				var err error
				req, err = http.NewRequest("GET", httpWebHomeURL, nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("is redirected to the same URL but with https", func() {
				resp, err := defaultTransport.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(http.StatusMovedPermanently))

				location, err := resp.Location()
				Expect(err).NotTo(HaveOccurred())

				Expect(location.Scheme).To(Equal("https"))
				Expect(location.Host).To(Equal(host))
				Expect(location.Path).To(Equal(path))
			})
		})

		Context("with forwarding header", func() {
			var client *http.Client
			var req *http.Request

			BeforeEach(func() {
				client = &http.Client{}
				req, _ = http.NewRequest("GET", httpWebHomeURL, nil)
				req.Header.Set("X-Forwarded-Proto", "https")
			})

			It("is accepted", func() {
				resp, err := client.Do(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})
})
