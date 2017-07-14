package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

const TIMEOUT_SECS = 45

//go:generate counterfeiter . HttpClient

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

func NewHttpClient(cfg config.Config) *http.Client {
	parsedUrl, _ := url.Parse(cfg.ApiURL)
	if parsedUrl.Scheme == "https" {
		return newHttpsClient(cfg)
	} else {
		return newHttpClient()
	}
}

func newHttpClient() *http.Client {
	return &http.Client{Timeout: time.Second * TIMEOUT_SECS}
}

func newHttpsClient(cfg config.Config) *http.Client {
	serverCaPath := cfg.CaCert
	trustedCAs := x509.NewCertPool()

	tlsConfig := &tls.Config{
		InsecureSkipVerify:       cfg.InsecureSkipVerify,
		PreferServerCipherSuites: true,
	}

	if len(serverCaPath) > 0 && !cfg.InsecureSkipVerify {
		for _, certPath := range serverCaPath {
			_, err := os.Stat(certPath)
			handleError(err)
			serverCA, err := ioutil.ReadFile(certPath)
			handleError(err)
			ok := trustedCAs.AppendCertsFromPEM([]byte(serverCA))

			if !ok {
				log.Fatal("failed to parse root certificate")
			}
		}

		tlsConfig.RootCAs = trustedCAs
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * TIMEOUT_SECS,
	}
	return client
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal("Fatal", err)
	}
}
