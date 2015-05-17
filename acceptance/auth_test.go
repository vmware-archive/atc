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

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

var _ = Describe("Auth", func() {
	var atcProcess ifrit.Process
	var dbListener *pq.Listener
	var atcPort uint16

	BeforeEach(func() {
		logger := lagertest.NewTestLogger("test")
		postgresRunner.CreateTestDB()
		dbConn = postgresRunner.Open()
		dbListener = pq.NewListener(postgresRunner.DataSourceName(), time.Second, time.Minute, nil)
		bus := db.NewNotificationsBus(dbListener)
		sqlDB = db.NewSQL(logger, dbConn, bus)

		Ω(sqlDB.SaveConfig(atc.DefaultPipelineName, atc.Config{}, db.ConfigVersion(1))).Should(Succeed())

		atcBin, err := gexec.Build("github.com/concourse/atc/cmd/atc")
		Ω(err).ShouldNot(HaveOccurred())

		atcPort = 5697 + uint16(GinkgoParallelNode())

		var atcCommand *exec.Cmd
		atcCommand, atcPort = createATCCommandWithFlags(
			atcBin,
			1,
			map[string]string{
				"-publiclyViewable": "false",
			})

		atcRunner := ginkgomon.New(ginkgomon.Config{
			Command:       atcCommand,
			Name:          "atc",
			StartCheck:    "atc.listening",
			AnsiColorCode: "32m",
		})
		atcProcess = ginkgomon.Invoke(atcRunner)
	})

	AfterEach(func() {
		ginkgomon.Interrupt(atcProcess)

		Ω(dbConn.Close()).Should(Succeed())
		Ω(dbListener.Close()).Should(Succeed())

		postgresRunner.DropTestDB()
	})

	It("can reach the page", func() {
		request, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:%d", atcPort), nil)

		resp, err := http.DefaultClient.Do(request)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(resp.StatusCode).Should(Equal(http.StatusUnauthorized))

		request.SetBasicAuth("admin", "password")
		resp, err = http.DefaultClient.Do(request)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(resp.StatusCode).Should(Equal(http.StatusOK))
	})
})
