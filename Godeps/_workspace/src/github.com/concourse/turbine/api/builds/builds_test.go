package builds_test

import (
	. "github.com/concourse/turbine/api/builds"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("merging", func() {
		It("merges params while preserving other properties", func() {
			Ω(Config{
				Image: "some-image",
				Params: map[string]string{
					"FOO": "1",
					"BAR": "2",
				},
			}.Merge(Config{
				Params: map[string]string{
					"FOO": "3",
					"BAZ": "4",
				},
			})).Should(Equal(Config{
				Image: "some-image",
				Params: map[string]string{
					"FOO": "3",
					"BAR": "2",
					"BAZ": "4",
				},
			}))
		})

		It("overrides the image", func() {
			Ω(Config{
				Image: "some-image",
				Run: RunConfig{
					Path: "some-path",
					Args: []string{"arg1", "arg2"},
				},
			}.Merge(Config{
				Image: "better-image",
			})).Should(Equal(Config{
				Image: "better-image",
				Run: RunConfig{
					Path: "some-path",
					Args: []string{"arg1", "arg2"},
				},
			}))
		})

		It("overrides the run config", func() {
			Ω(Config{
				Run: RunConfig{
					Path: "some-path",
					Args: []string{"arg1", "arg2"},
				},
			}.Merge(Config{
				Image: "some-image",
				Run: RunConfig{
					Path: "better-path",
					Args: []string{"better-arg1", "better-arg2"},
				},
			})).Should(Equal(Config{
				Image: "some-image",
				Run: RunConfig{
					Path: "better-path",
					Args: []string{"better-arg1", "better-arg2"},
				},
			}))
		})

		It("overrides the run config even with no args", func() {
			Ω(Config{
				Run: RunConfig{
					Path: "some-path",
					Args: []string{"arg1", "arg2"},
				},
			}.Merge(Config{
				Image: "some-image",
				Run: RunConfig{
					Path: "better-path",
				},
			})).Should(Equal(Config{
				Image: "some-image",
				Run: RunConfig{
					Path: "better-path",
				},
			}))
		})

		It("overrides input destinations", func() {
			Ω(Config{
				Inputs: []Input{
					{
						Name:            "some-input",
						DestinationPath: "some-destination",
					},
					{
						Name:            "another-input",
						DestinationPath: "another-destination",
					},
					{
						Name: "yet-another-input",
					},
				},
			}.Merge(Config{
				Inputs: []Input{
					{
						Name:            "another-input",
						DestinationPath: "an-even-better-destination",
					},
					{
						Name:            "yet-another-input",
						DestinationPath: "new-destination",
					},
				},
			})).Should(Equal(Config{
				Inputs: []Input{
					{
						Name:            "some-input",
						DestinationPath: "some-destination",
					},
					{
						Name:            "another-input",
						DestinationPath: "an-even-better-destination",
					},
					{
						Name:            "yet-another-input",
						DestinationPath: "new-destination",
					},
				},
			}))
		})
	})
})
