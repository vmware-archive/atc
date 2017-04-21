package genericoauth_test

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/concourse/atc/auth/genericoauth"
	"github.com/concourse/atc/auth/verifier"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthorityVerifier", func() {
	var verifier verifier.Verifier
	var httpClient *http.Client
	var verified bool
	var verifyErr error
	var jwtToken *jwt.Token

	BeforeEach(func() {

		verifier = NewAuthorityVerifier(
			"mainteam",
		)

		jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		})

		accessToken, err := jwtToken.SigningString()
		Expect(err).NotTo(HaveOccurred())

		oauthToken := &oauth2.Token{
			AccessToken: accessToken,
		}
		c := &oauth2.Config{}
		httpClient = c.Client(oauth2.NoContext, oauthToken)
	})

	JustBeforeEach(func() {
		verified, verifyErr = verifier.Verify(lagertest.NewTestLogger("test"), httpClient)
	})

	Context("when token does not contain 'authorities'", func() {
		It("user is not verified", func() {
			Expect(verified).To(BeFalse())
			Expect(verifyErr).To(HaveOccurred())
		})
	})

	Context("when user has proper authority", func() {
		BeforeEach(func() {
			jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp":       time.Now().Add(time.Hour * 72).Unix(),
				"authority": []string{"mainteam"},
			})

			accessToken, err := jwtToken.SigningString()
			Expect(err).NotTo(HaveOccurred())

			oauthToken := &oauth2.Token{
				AccessToken: accessToken,
			}
			c := &oauth2.Config{}
			httpClient = c.Client(oauth2.NoContext, oauthToken)
		})

		It("returns true", func() {
			Expect(verifyErr).NotTo(HaveOccurred())
			Expect(verified).To(BeTrue())
		})

	})

	Context("and user does not have proper authority", func() {
		BeforeEach(func() {
			jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp":       time.Now().Add(time.Hour * 72).Unix(),
				"authority": []string{"someotherteam"},
			})

			accessToken, err := jwtToken.SigningString()
			Expect(err).NotTo(HaveOccurred())

			oauthToken := &oauth2.Token{
				AccessToken: accessToken,
			}
			c := &oauth2.Config{}
			httpClient = c.Client(oauth2.NoContext, oauthToken)
		})

		It("returns false", func() {
			Expect(verifyErr).NotTo(HaveOccurred())
			Expect(verified).To(BeFalse())
		})
	})
})
