package acceptance_test

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/lib/pq"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/sclevine/agouti"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "github.com/sclevine/agouti/matchers"

	"github.com/cloudfoundry/gunk/urljoiner"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

var _ = Describe("Pipeline Pausing", func() {
	var atcProcess ifrit.Process
	var dbListener *pq.Listener
	var atcPort uint16
	var pipelineDBFactory db.PipelineDBFactory
	var pipelineDB db.PipelineDB
	var otherPipelineDB db.PipelineDB

	BeforeEach(func() {
		var err error
		atcBin, err := gexec.Build("github.com/concourse/atc/cmd/atc")
		Ω(err).ShouldNot(HaveOccurred())

		dbLogger := lagertest.NewTestLogger("test")
		postgresRunner.CreateTestDB()
		dbConn = postgresRunner.Open()
		dbListener = pq.NewListener(postgresRunner.DataSourceName(), time.Second, time.Minute, nil)
		bus := db.NewNotificationsBus(dbListener)
		sqlDB = db.NewSQL(dbLogger, dbConn, bus)
		pipelineDBFactory = db.NewPipelineDBFactory(dbLogger, dbConn, bus, sqlDB)

		var atcCommand *exec.Cmd
		atcCommand, atcPort = createATCCommand(atcBin, 1)
		atcProcess = startATC(atcCommand)
	})

	AfterEach(func() {
		ginkgomon.Interrupt(atcProcess)

		Ω(dbConn.Close()).Should(Succeed())
		Ω(dbListener.Close()).Should(Succeed())

		postgresRunner.DropTestDB()
	})

	Describe("pausing a pipeline", func() {
		var page *agouti.Page

		BeforeEach(func() {
			var err error
			page, err = agoutiDriver.NewPage()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(page.Destroy()).To(Succeed())
		})

		homepage := func() string {
			return fmt.Sprintf("http://127.0.0.1:%d/", atcPort)
		}

		withPath := func(path string) string {
			return urljoiner.Join(homepage(), path)
		}

		Context("with a job in the configuration", func() {

			BeforeEach(func() {
				var err error

				Ω(sqlDB.SaveConfig("some-pipeline", atc.Config{
					Jobs: []atc.JobConfig{
						{Name: "some-job-name"},
					},
				}, db.ConfigVersion(1))).Should(Succeed())

				Ω(sqlDB.SaveConfig("another-pipeline", atc.Config{
					Jobs: []atc.JobConfig{
						{Name: "another-job-name"},
					},
				}, db.ConfigVersion(1))).Should(Succeed())

				pipelineDB, err = pipelineDBFactory.BuildWithName("some-pipeline")
				Ω(err).ShouldNot(HaveOccurred())

				otherPipelineDB, err = pipelineDBFactory.BuildWithName("another-pipeline")
				Ω(err).ShouldNot(HaveOccurred())

			})

			homeLink := ".js-groups li:nth-of-type(2) a"
			defaultPipelineLink := ".js-pipelinesNav-list li:nth-of-type(1) a"
			anotherPipelineLink := ".js-pipelinesNav-list li:nth-of-type(2) a"
			anotherPipelineItem := ".js-pipelinesNav-list li:nth-of-type(2)"

			It("can pause the pipelines", func() {
				Expect(page.Navigate(homepage())).To(Succeed())
				// we will need to authenticate later to prove it is working for our page
				Authenticate(page, "admin", "password")

				Eventually(page.Find("#pipeline").Text).Should(ContainSubstring("some-job-name"))

				Expect(page.Find(".js-pipelinesNav-toggle").Click()).To(Succeed())

				Expect(page.Find(defaultPipelineLink)).To(HaveText("some-pipeline"))
				Expect(page.Find(anotherPipelineLink)).To(HaveText("another-pipeline"))

				Expect(page.Find(anotherPipelineLink).Click()).To(Succeed())

				Eventually(page).Should(HaveURL(withPath("/pipelines/another-pipeline")))
				Expect(page.Find(homeLink).Click()).To(Succeed())
				Eventually(page).Should(HaveURL(withPath("/pipelines/another-pipeline")))

				Expect(page.Find(".js-pipelinesNav-toggle").Click()).To(Succeed())
				Eventually(page.Find(defaultPipelineLink)).Should(HaveText("some-pipeline"))
				Eventually(page.Find("#pipeline").Text).Should(ContainSubstring("another-job-name"))

				Eventually(page.Find(anotherPipelineItem + " .js-pauseUnpause")).Should(BeVisible())
				Eventually(page.Find(anotherPipelineItem + " .js-pauseUnpause.disabled")).Should(BeFound())

				Expect(page.Find(anotherPipelineItem + " .js-pauseUnpause").Click()).To(Succeed())
				Eventually(page.Find(anotherPipelineItem + " .js-pauseUnpause.enabled")).Should(BeFound())

				// top bar should show the pipeline is paused
				Eventually(page.Find(".js-groups.paused")).Should(BeFound())

				page.Refresh()

				Eventually(page.Find(".js-groups.paused")).Should(BeFound())
				Expect(page.Find(".js-pipelinesNav-toggle").Click()).To(Succeed())
				Eventually(page.Find(anotherPipelineItem + " .js-pauseUnpause")).Should(BeVisible())
				Eventually(page.Find(anotherPipelineItem + " .js-pauseUnpause.enabled")).Should(BeFound())

				Expect(page.Find(anotherPipelineItem + " .js-pauseUnpause").Click()).To(Succeed())
				Eventually(page.Find(anotherPipelineItem + " .js-pauseUnpause.disabled")).Should(BeFound())

				Consistently(page.Find(".js-groups.paused")).ShouldNot(BeFound())

				page.Refresh()

				Expect(page.Find(".js-pipelinesNav-toggle").Click()).To(Succeed())
				Eventually(page.Find(anotherPipelineItem + " .js-pauseUnpause")).Should(BeVisible())
				Eventually(page.Find(anotherPipelineItem + " .js-pauseUnpause.disabled")).Should(BeFound())
			})
		})
	})
})
