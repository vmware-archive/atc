package credhub

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/permissions"

// GetPermissions returns the permissions of a credential.
func (ch *CredHub) GetPermissions(credName string) ([]permissions.Permission, error) {
	panic("Not implemented")
}

// AddPermissions adds permissions to a credential.
func (ch *CredHub) AddPermissions(credName string, perms []permissions.Permission) ([]permissions.Permission, error) {
	panic("Not implemented")
}

// DeletePermissions deletes permissions on a credential by actor.
func (ch *CredHub) DeletePermissions(credName string, actor string) error {
	panic("Not implemented")
}
