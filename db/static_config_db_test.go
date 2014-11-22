package db_test

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StaticConfigDB", func() {
	var staticConfigDB db.StaticConfigDB
	var config atc.Config

	var public bool
	var err error

	JustBeforeEach(func() {
		staticConfigDB = db.StaticConfigDB{config}

		public, err = staticConfigDB.JobIsPublic("some-job")
	})

	Describe("determining if a job's builds are publically viewable", func() {
		Context("when the job is publically viewable", func() {
			BeforeEach(func() {
				config = atc.Config{
					Jobs: atc.JobConfigs{
						{
							Name:   "some-job",
							Public: true,
						},
					},
				}
			})

			It("returns true", func() {
				Ω(public).Should(BeTrue())
			})

			It("does not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when the job is not publically viewable", func() {
			BeforeEach(func() {
				config = atc.Config{
					Jobs: atc.JobConfigs{
						{
							Name:   "some-job",
							Public: false,
						},
					},
				}
			})

			It("returns false", func() {
				Ω(public).Should(BeFalse())
			})

			It("does not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when the job with the given name can't be found", func() {
			BeforeEach(func() {
				config = atc.Config{
					Jobs: atc.JobConfigs{
						{
							Name:   "DIFFERENT-JOB",
							Public: false,
						},
					},
				}
			})

			It("errors", func() {
				Ω(err).Should(HaveOccurred())
			})
		})
	})
})
