/*
Package credhub is a client library for interacting with a CredHub server.

More information on CredHub can be found at https://github.com/cloudfoundry-incubator/credhub

Server HTTP API specification can be found at http://credhub-api.cfapps.io
*/
package credhub

import (
	"net/http"
	"net/url"

	"crypto/x509"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

// CredHub client to access CredHub APIs.
//
// Use New() to construct a new CredHub object, which can then interact with the CredHub API.
type CredHub struct {
	// ApiURL is the host and port of the CredHub server to target
	// Example: https://credhub.example.com:8844
	ApiURL string

	// Auth provides an authentication Strategy for authenticated requests to the CredHub server
	// Can be type asserted to a specific Strategy type to get additional functionality and information.
	// eg. auth.OAuthStrategy provides Logout(), Refresh(), AccessToken() and RefreshToken()
	Auth auth.Strategy

	baseURL       *url.URL
	defaultClient *http.Client

	// Trusted CA certificates in PEM format for making TLS connections to CredHub and auth servers
	caCerts *x509.CertPool

	// Skip certificate verification of TLS connections to CredHub and auth servers. Not recommended!
	insecureSkipVerify bool

	authBuilder auth.Builder
	authURL     *url.URL
}
