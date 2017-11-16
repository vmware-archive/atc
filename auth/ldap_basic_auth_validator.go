package auth

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/db"
	ldap "gopkg.in/ldap.v2"
)

type ldapBasicAuthValidator struct {
	team db.Team
}

type groupTemplateVars struct {
	UserDN string
}

// NewLdapBasicAuthValidator creates an LDAP validator
func NewLdapBasicAuthValidator(team db.Team) Validator {
	return ldapBasicAuthValidator{
		team: team,
	}
}

func (v ldapBasicAuthValidator) IsAuthenticated(logger lager.Logger, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	username, password, err := extractUsernameAndPassword(auth)
	if err != nil {
		return false
	}

	return v.correctCredentials(
		logger, username, password,
	)
}

func (v ldapBasicAuthValidator) correctCredentials(
	logger lager.Logger, checkUsername string, checkPassword string,
) bool {

	ldapAuthConfig := v.team.LdapBasicAuth()

	if ldapAuthConfig == nil {
		logger.Error("ldapAuthConfig is nil. Skipping LDAP Auth check", errors.New("ldapAuthConfig is nil"))
		return false
	}

	rootCAs := x509.NewCertPool()
	ldapConn := &ldap.Conn{}
	var err error

	if ldapAuthConfig.TLSEnabled {
		if ldapAuthConfig.TLSCA == "" {
			rootCAs, err = x509.SystemCertPool()
			if err != nil {
				logger.Error("Could not load system certificate store: ", err)
				return false
			}
		} else {
			rootCAs = x509.NewCertPool()
			rootCAs.AppendCertsFromPEM([]byte(ldapAuthConfig.TLSCA))
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: ldapAuthConfig.TLSInsecureSkipVerify,
			RootCAs:            rootCAs,
			ServerName:         ldapAuthConfig.Server,
		}
		ldapConn, err = ldap.DialTLS("tcp",
			fmt.Sprintf("%s:%d", ldapAuthConfig.Server, ldapAuthConfig.Port),
			tlsConfig)
		if err != nil {
			logger.Error("Error connecting to LDAP server: ", err)
			return false
		}
	} else {
		ldapConn, err = ldap.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", ldapAuthConfig.Server, ldapAuthConfig.Port),
		)
		if err != nil {
			logger.Error("Error connecting to LDAP server", err)
			return false
		}
	}

	defer ldapConn.Close()

	// First bind with a read only user
	err = ldapConn.Bind(
		ldapAuthConfig.BindUsername,
		ldapAuthConfig.BindPassword)
	if err != nil {
		logger.Error("Error binding to LDAP server: ", err)
		return false
	}

	userFilter := fmt.Sprintf("(&(objectClass=user)(%s=%s))",
		ldap.EscapeFilter(ldapAuthConfig.UserAttribute),
		ldap.EscapeFilter(checkUsername))

	// Search for the given username and retrieve the DN
	searchRequest := ldap.NewSearchRequest(
		ldapAuthConfig.UserBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		userFilter,
		[]string{"dn"},
		nil,
	)

	sr, err := ldapConn.Search(searchRequest)
	if err != nil {
		logger.Error("Error executing LDAP search for user: ", err)
		return false
	}

	if len(sr.Entries) != 1 {
		logger.Info("User not found in LDAP")
		return false
	}

	userdn := sr.Entries[0].DN

	groupFilterTemplate, err := template.New("groupFilter").Parse(ldapAuthConfig.GroupFilter)
	if err != nil {
		logger.Error("Failed to parse GroupFilter template: ", err)
		return false
	}
	groupFilterBuffer := new(bytes.Buffer)
	templateVars := groupTemplateVars{
		UserDN: ldap.EscapeFilter(userdn),
	}
	err = groupFilterTemplate.Execute(groupFilterBuffer, templateVars)
	if err != nil {
		logger.Error("Failed to render GroupFilter template: ", err)
		return false
	}
	groupFilter := groupFilterBuffer.String()

	logger.Debug(groupFilter)

	// Enumerate the user's group membership
	groupSearchRequest := ldap.NewSearchRequest(
		ldapAuthConfig.GroupBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		groupFilter,
		[]string{"dn"},
		nil,
	)

	gsr, err := ldapConn.Search(groupSearchRequest)
	if err != nil {
		logger.Error("Error occurred searching LDAP for team GroupDN: ", err)
		return false
	}

	userInGroup := false
	for _, entry := range gsr.Entries {
		logger.Debug("Group membership found: " + entry.DN)
		if entry.DN == ldapAuthConfig.GroupDN {
			userInGroup = true
			break
		}
	}
	if !userInGroup {
		logger.Info("User is not a member of the team's group")
		return false
	}
	// Bind as the user to verify their password
	err = ldapConn.Bind(userdn, checkPassword)
	if err != nil {
		logger.Error("Authentication error: invalid username or password: ", err)
		return false
	} else {
		return true
	}

	return false
}
