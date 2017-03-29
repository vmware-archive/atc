package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"time"

	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/auth/authfakes"
	"github.com/concourse/atc/db/dbfakes"
)

var _ = Describe("LogOutHandler", func() {
	Describe("GET /auth/logout", func() {
		var (
			fakeProviderFactory *authfakes.FakeProviderFactory
			signingKey          *rsa.PrivateKey
			server              *httptest.Server
			client              *http.Client
			request             *http.Request
			response            *http.Response
			err                 error
			expire              time.Duration
			httpOnly            bool
			secure              bool
		)

		BeforeEach(func() {
			fakeProviderFactory = new(authfakes.FakeProviderFactory)
			fakeTeamDBFactory := new(dbfakes.FakeTeamDBFactory)
			signingKey, err = rsa.GenerateKey(rand.Reader, 1024)
			Expect(err).ToNot(HaveOccurred())
			expire = 24 * time.Hour
			httpOnly = true
			secure = false

			handler, err := auth.NewOAuthHandler(
				lagertest.NewTestLogger("test"),
				fakeProviderFactory,
				fakeTeamDBFactory,
				signingKey,
				expire,
				httpOnly,
				secure,
			)
			Expect(err).ToNot(HaveOccurred())

			mux := http.NewServeMux()
			mux.Handle("/auth/", handler)

			server = httptest.NewServer(mux)

			client = &http.Client{
				Transport: &http.Transport{},
			}

			request, err = http.NewRequest("GET", server.URL+"/auth/logout", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		It("deletes ATC-Authorization cookie", func() {
			cookies := response.Cookies()
			Expect(len(cookies)).To(Equal(1))

			deletedCookie := cookies[0]
			Expect(deletedCookie.Name).To(Equal(auth.CookieName))
			Expect(deletedCookie.MaxAge).To(Equal(-1))
		})
	})
})
