package credhub

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"

// Regenerate generates and returns a new credential version using the same parameters existing credential. The returned credential may be of any type.
func (ch *CredHub) Regenerate(name string) (credentials.Credential, error) {
	panic("Not implemented")
}

// RegeneratePassword generates and returns a new credential version using the same parameters existing credential. The returned credential must be of type 'password'.
func (ch *CredHub) RegeneratePassword(name string) (credentials.Password, error) {
	panic("Not implemented")
}

// RegenerateUser generates and returns a new credential version using the same parameters existing credential. The returned credential must be of type 'user'.
func (ch *CredHub) RegenerateUser(name string) (credentials.User, error) {
	panic("Not implemented")
}

// RegenerateCertificate generates and returns a new credential version using the same parameters existing credential. The returned credential must be of type 'certificate'.
func (ch *CredHub) RegenerateCertificate(name string) (credentials.Certificate, error) {
	panic("Not implemented")
}

// RegenerateRSA generates and returns a new credential version using the same parameters existing credential. The returned credential must be of type 'rsa'.
func (ch *CredHub) RegenerateRSA(name string) (credentials.RSA, error) {
	panic("Not implemented")
}

// RegenerateSSH generates and returns a new credential version using the same parameters existing credential. The returned credential must be of type 'ssh'.
func (ch *CredHub) RegenerateSSH(name string) (credentials.SSH, error) {
	panic("Not implemented")
}
