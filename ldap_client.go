package atc

import (
	"crypto/tls"

	ldap "gopkg.in/ldap.v2"
)

type LdapDialer interface {
	Dial(network, addr string) (*ldap.Conn, error)
	DialTLS(network, addr string, config *tls.Config) (*ldap.Conn, error)
}

type LdapDialerImpl struct {
}

func (l *LdapDialerImpl) Dial(network, addr string) (*ldap.Conn, error) {
	return ldap.Dial(network, addr)
}

func (l *LdapDialerImpl) DialTLS(network, addr string, config *tls.Config) (*ldap.Conn, error) {
	return ldap.DialTLS(network, addr, config)
}
