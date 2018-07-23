package migrator_test

import (
	. "github.com/concourse/atc/db/migration/migrator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Context("SQL migrations", func() {
		It("parses the migration into statements", func() {
			migrationFileNames := make([]string, 1)
			migrationFileNames[0] = "1513895878_update_timestamp_with_timezone.up.sql"

			parser := NewParser()
			statements, err := parser.ParseFile(migrationFileNames)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(statements)).To(Equal(29))
		})
	})
})
