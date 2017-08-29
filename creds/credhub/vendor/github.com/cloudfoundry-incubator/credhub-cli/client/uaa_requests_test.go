package client_test

import (
	"net/http"

	. "github.com/cloudfoundry-incubator/credhub-cli/client"

	"bytes"

	b64 "encoding/base64"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA Requests", func() {
	var cfg config.Config

	BeforeEach(func() {
		cfg = config.Config{
			AuthURL: "http://example.com",
		}
	})

	Describe("NewPasswordGrantTokenRequest", func() {
		It("Returns a request for the uaa oauth token endpoint", func() {
			user := "my-user"
			pass := "my-pass"

			basicEncoded := b64.StdEncoding.EncodeToString([]byte(config.AuthClient + ":"))

			request := NewPasswordGrantTokenRequest(cfg, user, pass)
			Expect(request.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
			Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/x-www-form-urlencoded"}))
			Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Basic " + basicEncoded}))
			Expect(request.URL.Path).To(Equal("/oauth/token/"))
			Expect(request.Method).To(Equal("POST"))

			byteBuff := new(bytes.Buffer)
			byteBuff.ReadFrom(request.Body)

			Expect(byteBuff.String()).To(ContainSubstring("grant_type=password"))
			Expect(byteBuff.String()).To(ContainSubstring("password=my-pass"))
			Expect(byteBuff.String()).To(ContainSubstring("response_type=token"))
			Expect(byteBuff.String()).To(ContainSubstring("username=my-user"))
		})
	})

	Describe("NewClientCredentialsGrantTokenRequest", func() {
		It("Returns a request for the uaa oauth token endpoint", func() {
			client := "my-client"
			clientSecret := "my-client-secret"

			request := NewClientCredentialsGrantTokenRequest(cfg, client, clientSecret)
			Expect(request.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
			Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/x-www-form-urlencoded"}))
			Expect(request.URL.Path).To(Equal("/oauth/token/"))
			Expect(request.Method).To(Equal("POST"))

			byteBuff := new(bytes.Buffer)
			byteBuff.ReadFrom(request.Body)

			Expect(byteBuff.String()).To(ContainSubstring("grant_type=client_credentials"))
			Expect(byteBuff.String()).To(ContainSubstring("response_type=token"))
			Expect(byteBuff.String()).To(ContainSubstring("client_id=my-client"))
			Expect(byteBuff.String()).To(ContainSubstring("client_secret=my-client-secret"))
		})
	})

	Describe("NewRefreshTokenRequest", func() {
		It("Returns a request for the uaa oauth token endpoint to get refresh token", func() {
			request := NewRefreshTokenRequest(cfg)

			basicEncoded := b64.StdEncoding.EncodeToString([]byte(config.AuthClient + ":"))

			Expect(request.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
			Expect(request.Header).To(HaveKeyWithValue("Content-Type", []string{"application/x-www-form-urlencoded"}))
			Expect(request.Header).To(HaveKeyWithValue("Authorization", []string{"Basic " + basicEncoded}))
			Expect(request.URL.Path).To(Equal("/oauth/token/"))
			Expect(request.Method).To(Equal("POST"))

			byteBuff := new(bytes.Buffer)
			byteBuff.ReadFrom(request.Body)

			Expect(byteBuff.String()).To(ContainSubstring("grant_type=refresh_token"))
			Expect(byteBuff.String()).To(ContainSubstring("refresh_token=" + cfg.RefreshToken))
		})
	})

	Describe("NewTokenRevocationRequest", func() {
		It("Returns a request to revoke a refresh token", func() {
			cfg.RefreshToken = "5b9c9fd51ba14838ac2e6b222d487106-r"
			cfg.AccessToken = "defgh"
			expectedRequest, _ := http.NewRequest(
				"DELETE",
				cfg.AuthURL+"/oauth/token/revoke/5b9c9fd51ba14838ac2e6b222d487106-r",
				nil)
			expectedRequest.Header.Add("Authorization", "Bearer defgh")

			request, _ := NewTokenRevocationRequest(cfg)

			Expect(request).To(Equal(expectedRequest))
		})
	})
})
