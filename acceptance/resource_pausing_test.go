package acceptance_test

import (
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/sclevine/agouti"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"

	"code.cloudfoundry.org/gunk/urljoiner"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

var _ = Describe("Resource Pausing", func() {
	var atcCommand *ATCCommand
	var dbListener *pq.Listener
	var pipelineDB db.PipelineDB

	BeforeEach(func() {
		postgresRunner.Truncate()
		dbConn = db.Wrap(postgresRunner.Open())
		dbListener = pq.NewListener(postgresRunner.DataSourceName(), time.Second, time.Minute, nil)
		bus := db.NewNotificationsBus(dbListener, dbConn)

		sqlDB = db.NewSQL(dbConn, bus)

		atcCommand = NewATCCommand(atcBin, 1, postgresRunner.DataSourceName(), []string{}, BASIC_AUTH)
		err := atcCommand.Start()
		Expect(err).NotTo(HaveOccurred())

		teamDBFactory := db.NewTeamDBFactory(dbConn, bus)
		teamDB := teamDBFactory.GetTeamDB(atc.DefaultTeamName)
		// job build data
		_, _, err = teamDB.SaveConfig("some-pipeline", atc.Config{
			Jobs: atc.JobConfigs{
				{
					Name: "job-name",
					Plan: atc.PlanSequence{
						{
							Get: "resource-name",
						},
					},
				},
			},
			Resources: atc.ResourceConfigs{
				{Name: "resource-name"},
			},
		}, db.ConfigVersion(1), db.PipelineUnpaused)
		Expect(err).NotTo(HaveOccurred())

		savedPipeline, found, err := teamDB.GetPipelineByName("some-pipeline")
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue())

		pipelineDBFactory := db.NewPipelineDBFactory(dbConn, bus)
		pipelineDB = pipelineDBFactory.Build(savedPipeline)
	})

	AfterEach(func() {
		atcCommand.Stop()

		Expect(dbConn.Close()).To(Succeed())
		Expect(dbListener.Close()).To(Succeed())
	})

	Describe("pausing a resource", func() {
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
			return atcCommand.URL("")
		}

		withPath := func(path string) string {
			return urljoiner.Join(homepage(), path)
		}

		It("can view the resource", func() {
			// homepage -> resource detail
			Login(page, homepage())

			Expect(page.Navigate(homepage())).To(Succeed())
			Eventually(page.FindByLink("resource-name")).Should(BeFound())
			Expect(page.FindByLink("resource-name").Click()).To(Succeed())

			// resource detail -> paused resource detail
			Eventually(page).Should(HaveURL(withPath("/teams/main/pipelines/some-pipeline/resources/resource-name")))
			Expect(page.Find("h1")).To(HaveText("resource-name"))

			Expect(page.Find(".js-resource .js-pauseUnpause").Click()).To(Succeed())

			Expect(page.Navigate(homepage())).To(Succeed())
			Eventually(page.FindByLink("resource-name")).Should(BeFound())
			Expect(page.FindByLink("resource-name").Click()).To(Succeed())
			Expect(page.Find(".js-resource .js-pauseUnpause").Click()).To(Succeed())

			Eventually(page.Find(".header i.fa-play")).Should(BeFound())

			resource, _, err := pipelineDB.GetResource("resource-name")
			Expect(err).NotTo(HaveOccurred())

			err = pipelineDB.SetResourceCheckError(resource, errors.New("failed to foo the bar"))
			Expect(err).NotTo(HaveOccurred())

			page.Refresh()

			Eventually(page.Find(".header h3")).Should(HaveText("checking failed"))
			Eventually(page.Find(".build-step .step-body")).Should(HaveText("failed to foo the bar"))
		})
	})
})
