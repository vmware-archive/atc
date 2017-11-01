package atc

import (
	"errors"

	multierror "github.com/hashicorp/go-multierror"
)

type AuthFlags struct {
	NoAuth bool `long:"no-really-i-dont-want-any-auth" description:"Ignore warnings about not configuring auth"`

	BasicAuth     BasicAuthFlag     `group:"Basic Authentication" namespace:"basic-auth"`
	LdapBasicAuth LdapBasicAuthFlag `group:"LDAP Authentication" namespace:"ldap-auth"`
}

type BasicAuthFlag struct {
	Username string `long:"username" description:"Username to use for basic auth."`
	Password string `long:"password" description:"Password to use for basic auth."`
}

func (auth *BasicAuthFlag) IsConfigured() bool {
	return auth.Username != "" || auth.Password != ""
}

func (auth *BasicAuthFlag) Validate() error {
	var errs *multierror.Error
	if auth.Username == "" {
		errs = multierror.Append(
			errs,
			errors.New("must specify --basic-auth-username to use basic auth."),
		)
	}
	if auth.Password == "" {
		errs = multierror.Append(
			errs,
			errors.New("must specify --basic-auth-password to use basic auth."),
		)
	}
	return errs.ErrorOrNil()
}

type LdapBasicAuthFlag struct {
	Server                string `long:"server" description:"Server to use for LDAP auth."`
	Port                  uint16 `long:"port" description:"TCP Port number to use for LDAP auth."`
	TLSEnabled            bool   `long:"tls-enabled" description:"Connect to LDAP with TLS."`
	TLSInsecureSkipVerify bool   `long:"tls-insecure-skip-verify" description:"Skip LDAP server certificate validation when connecting via TLS."`
	TLSCA                 string `long:"tls-ca" description:"CA Certificate to verify the server certificate when connecting to LDAP with TLS."`
	BindUsername          string `long:"bind-username" description:"Username to use to validate LDAP users."`
	BindPassword          string `long:"bind-password" description:"Password to use to validate LDAP users."`
	UserBaseDN            string `long:"user-base-dn" description:"Base DN to use when searching for LDAP user."`
	GroupBaseDN           string `long:"group-base-dn" description:"Base DN to use when searching LDAP for the user's groups."`
	GroupDN               string `long:"group-dn" description:"DN of Group to authorize access to the team."`
	UserAttribute         string `long:"user-attribute" description:"Unique LDAP Attribute that contains username."`
	GroupFilter           string `long:"group-filter" description:"LDAP Filter to use when querying user's group membership."`
}

func (auth *LdapBasicAuthFlag) IsConfigured() bool {
	return auth.Server != ""
}

func (auth *LdapBasicAuthFlag) Validate() error {
	var errs *multierror.Error
	if auth.Server == "" {
		errs = multierror.Append(
			errs,
			errors.New("must specify --ldap-auth-server to use LDAP auth."),
		)
	}
	if auth.Port == 0 {
		errs = multierror.Append(
			errs,
			errors.New("must specify --ldap-auth-port to use LDAP auth."),
		)
	}
	if auth.BindUsername == "" {
		errs = multierror.Append(
			errs,
			errors.New("must specify --ldap-auth-bind-username to use LDAP auth."),
		)
	}
	return errs.ErrorOrNil()
}
