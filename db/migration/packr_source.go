package migration

import "github.com/gobuffalo/packr"

type PackrSource struct {
	packr.Box
}

func (s *PackrSource) AssetNames() []string {
	migrations := []string{}
	for _, name := range s.Box.List() {
		if name != "migrations.go" {
			migrations = append(migrations, name)
		}
	}

	return migrations
}

func (s *PackrSource) Asset(name string) ([]byte, error) {
	return s.Box.MustBytes(name)
}
