package acceptance_test

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"

	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/postgresrunner"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"

	"testing"
	"time"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var postgresRunner postgresrunner.Runner
var dbConn *sql.DB
var dbProcess ifrit.Process

var sqlDB *db.SQLDB

var agoutiDriver *agouti.WebDriver

var _ = BeforeSuite(func() {
	postgresRunner = postgresrunner.Runner{
		Port: 5432 + GinkgoParallelNode(),
	}

	dbProcess = ifrit.Envoke(postgresRunner)

	agoutiDriver = agouti.PhantomJS()
	Expect(agoutiDriver.Start()).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(agoutiDriver.Stop()).To(Succeed())

	dbProcess.Signal(os.Interrupt)
	Eventually(dbProcess.Wait(), 10*time.Second).Should(Receive())
})

func Screenshot(page *agouti.Page) {
	page.Screenshot("/tmp/screenshot.png")
}

func Authenticate(page *agouti.Page, username, password string) {
	header := fmt.Sprintf("%s:%s", username, password)

	page.SetCookie(&http.Cookie{
		Name:  auth.CookieName,
		Value: "Basic " + base64.StdEncoding.EncodeToString([]byte(header)),
	})

	// PhantomJS won't send the cookie on ajax requests if the page is not
	// refreshed
	page.Refresh()
}

func createATCCommandWithFlags(atcBin string, atcServerNumber uint16, flags map[string]string) (*exec.Cmd, uint16) {
	atcPort := 5697 + uint16(GinkgoParallelNode()) + (atcServerNumber * 100)
	debugPort := 6697 + uint16(GinkgoParallelNode()) + (atcServerNumber * 100)

	atcFlags := map[string]string{
		"-publiclyViewable":   "true",
		"-forceHTTPS":         "false",
		"-webListenPort":      fmt.Sprintf("%d", atcPort),
		"-callbacksURL":       fmt.Sprintf("http://127.0.0.1:%d", atcPort),
		"-debugListenPort":    fmt.Sprintf("%d", debugPort),
		"-httpUsername":       "admin",
		"-httpHashedPassword": "$2a$04$DYaOWeQgyxTCv7QxydTP9u1KnwXWSKipC4BeTuBy.9m.IlkAdqNGG", // "password"
		"-templates":          filepath.Join("..", "web", "templates"),
		"-public":             filepath.Join("..", "web", "public"),
		"-sqlDataSource":      postgresRunner.DataSourceName(),
	}

	// merge in any provided values
	if flags != nil {
		for k, _ := range atcFlags {
			if flags[k] != "" {
				atcFlags[k] = flags[k]
			}
		}
	}

	var flagsAsStrings []string
	for k, v := range atcFlags {
		// boolean flag values are treated differently;
		// they require = when setting to false.
		// See https://golang.org/pkg/flag/
		if v == "true" || v == "false" {
			flag := fmt.Sprintf("%s=%s", k, v)
			flagsAsStrings = append(flagsAsStrings, flag)
		} else {
			flagsAsStrings = append(flagsAsStrings, k, v)
		}

	}

	return exec.Command(
		atcBin,
		flagsAsStrings...,
	), atcPort
}

func createATCCommand(atcBin string, atcServerNumber uint16) (*exec.Cmd, uint16) {
	return createATCCommandWithFlags(atcBin, atcServerNumber, nil)
}

func startATC(atcCommand *exec.Cmd) ifrit.Process {
	atcRunner := ginkgomon.New(ginkgomon.Config{
		Command:       atcCommand,
		Name:          "atc",
		StartCheck:    "atc.listening",
		AnsiColorCode: "32m",
	})

	return ginkgomon.Invoke(atcRunner)
}
