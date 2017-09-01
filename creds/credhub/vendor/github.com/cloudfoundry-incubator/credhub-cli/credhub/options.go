package credhub

import (
	"crypto/x509"
	"errors"
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

// Option can be provided to New() to specify additional parameters for
// connecting to the CredHub server
type Option func(*CredHub) error

// Auth specifies the authentication Strategy. See the auth package
// for a full list of supported strategies.
func Auth(method auth.Builder) Option {
	return func(c *CredHub) error {
		c.authBuilder = method
		return nil
	}
}

// AuthURL specifies the authentication server for the OAuth strategy.
// If AuthURL provided, the AuthURL will be fetched from /info.
func AuthURL(authURL string) Option {
	return func(c *CredHub) error {
		var err error
		c.authURL, err = url.Parse(authURL)
		return err
	}
}

// CaCerts specifies the root certificates for HTTPS connections with the CredHub server.
//
// If the OAuthStrategy is used for Auth, the root certificates will also be used for HTTPS
// connections with the OAuth server.
func CaCerts(certs ...string) Option {
	return func(c *CredHub) error {
		c.caCerts = x509.NewCertPool()

		for _, cert := range certs {
			ok := c.caCerts.AppendCertsFromPEM([]byte(cert))
			if !ok {
				return errors.New("provided ca certs are invalid")
			}
		}

		return nil
	}
}

// SkipTLSValidation will skip root certificate verification for HTTPS. Not recommended!
func SkipTLSValidation() Option {
	return func(c *CredHub) error {
		c.insecureSkipVerify = true
		return nil
	}
}
