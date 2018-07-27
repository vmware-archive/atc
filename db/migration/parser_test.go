package migration_test

import (
	. "github.com/concourse/atc/db/migration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Context("SQL migrations", func() {
		It("parses the migration into statements", func() {
			migrationFileName := "1513895878_update_timestamp_with_timezone.up.sql"

			parser := NewParser()
			statements, err := parser.ParseFile(migrationFileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(statements)).To(Equal(13))
		})

		It("combines sql functions in one statement", func() {
			migrationFileName := "1530823998_create_teams_trigger.up.sql"

			parser := NewParser()
			statements, err := parser.ParseFile(migrationFileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(statements)).To(Equal(6))

		})

		It("removes the BEGIN and COMMIT statements", func() {
			migrationFileName := "1510670987_update_unique_constraint_for_resource_caches.down.sql"

			parser := NewParser()
			statements, err := parser.ParseFile(migrationFileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(statements)).To(Equal(2))
			Expect(statements[0]).ToNot(Equal("BEGIN"))
		})

		Context("No transactions", func() {
			It("marks migration as no transaction", func() {
			})
		})
	})

})
