package db_test

import (
	"github.com/concourse/atc"
	. "github.com/concourse/atc/db"
	"github.com/concourse/atc/db/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigDBWithDefaults", func() {
	var realConfigDB *fakes.FakeConfigDB
	var configDB ConfigDB
	var config atc.Config

	JustBeforeEach(func() {
		realConfigDB.GetConfigReturns(config, nil)
	})

	BeforeEach(func() {
		realConfigDB = new(fakes.FakeConfigDB)

		configDB = ConfigDBWithDefaults{
			ConfigDB: realConfigDB,
		}
	})

	Context("when an input does not specify its name or whether to trigger", func() {
		BeforeEach(func() {
			config = atc.Config{
				Jobs: atc.JobConfigs{
					{
						Name: "some-job",
						Inputs: []atc.InputConfig{
							{
								Resource: "foo",
							},
						},
					},
				},
			}
		})

		It("defaults trigger to true, and the name to the resource", func() {
			triggerDefault := true

			Ω(configDB.GetConfig()).Should(Equal(atc.Config{
				Jobs: atc.JobConfigs{
					{
						Name: "some-job",
						Inputs: []atc.InputConfig{
							{
								Name:     "foo",
								Resource: "foo",
								Trigger:  &triggerDefault,
							},
						},
					},
				},
			}))
		})
	})

	Context("when an output does not specify when to perform", func() {
		BeforeEach(func() {
			config = atc.Config{
				Jobs: atc.JobConfigs{
					{
						Name: "some-job",
						Outputs: []atc.OutputConfig{
							{
								Resource: "foo",
							},
						},
					},
				},
			}
		})

		It("performs on success", func() {
			Ω(configDB.GetConfig()).Should(Equal(atc.Config{
				Jobs: atc.JobConfigs{
					{
						Name: "some-job",
						Outputs: []atc.OutputConfig{
							{
								Resource:  "foo",
								PerformOn: []atc.OutputCondition{"success"},
							},
						},
					},
				},
			}))
		})
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
				public, _ := configDB.JobIsPublic("some-job")
				Ω(public).Should(BeTrue())
			})

			It("does not error", func() {
				_, err := configDB.JobIsPublic("some-job")
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
				public, _ := configDB.JobIsPublic("some-job")
				Ω(public).Should(BeFalse())
			})

			It("does not error", func() {
				_, err := configDB.JobIsPublic("some-job")
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when the job with the given name can't be found", func() {
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

			It("errors", func() {
				_, err := configDB.JobIsPublic("does-not-exist")
				Ω(err).Should(HaveOccurred())
			})
		})
	})
})
