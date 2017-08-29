package credhub

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
)

// Client provides an unauthenticated http.Client to the CredHub server
func (ch *CredHub) Client() *http.Client {
	if ch.defaultClient == nil {
		ch.defaultClient = ch.client()
	}

	return ch.defaultClient
}

func (ch *CredHub) client() *http.Client {
	if ch.baseURL.Scheme == "https" {
		return httpsClient(ch.insecureSkipVerify, ch.caCerts)
	} else {
		return httpClient()
	}
}

func httpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 45,
	}
}

func httpsClient(insecureSkipVerify bool, rootCAs *x509.CertPool) *http.Client {
	client := httpClient()

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify:       insecureSkipVerify,
			PreferServerCipherSuites: true,
			RootCAs:                  rootCAs,
		},
	}

	return client
}
