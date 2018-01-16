package inputmapper_test

import (
	"errors"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/db/dbfakes"
	"github.com/concourse/atc/scheduler/inputmapper"
	"github.com/concourse/atc/scheduler/inputmapper/inputconfig/inputconfigfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inputmapper", func() {
	var (
		fakePipeline    *dbfakes.FakePipeline
		fakeTransformer *inputconfigfakes.FakeTransformer

		inputMapper inputmapper.InputMapper

		disaster error
	)

	BeforeEach(func() {
		fakePipeline = new(dbfakes.FakePipeline)
		fakeTransformer = new(inputconfigfakes.FakeTransformer)

		inputMapper = inputmapper.NewInputMapper(fakePipeline, fakeTransformer)

		disaster = errors.New("bad thing")
	})

	Describe("SaveNextInputMapping", func() {
		var (
			versionsDB         *algorithm.VersionsDB
			fakeJob            *dbfakes.FakeJob
			fakeJobCombination *dbfakes.FakeJobCombination
			inputMapping       algorithm.InputMapping
			mappingErr         error
		)

		BeforeEach(func() {
			versionsDB = &algorithm.VersionsDB{
				JobCombinationIDs: map[string]int{"some-job": 1, "upstream": 2},
				ResourceIDs:       map[string]int{"a": 11, "b": 12, "no-versions": 13},
				ResourceVersions: []algorithm.ResourceVersion{
					{VersionID: 1, ResourceID: 11, CheckOrder: 1},
					{VersionID: 2, ResourceID: 12, CheckOrder: 1},
				},
				BuildOutputs: []algorithm.BuildOutput{
					{
						ResourceVersion:  algorithm.ResourceVersion{VersionID: 1, ResourceID: 11, CheckOrder: 1},
						BuildID:          98,
						JobCombinationID: 2,
					},
					{
						ResourceVersion:  algorithm.ResourceVersion{VersionID: 2, ResourceID: 12, CheckOrder: 1},
						BuildID:          99,
						JobCombinationID: 2,
					},
				},
			}
		})

		JustBeforeEach(func() {
			inputMapping, mappingErr = inputMapper.SaveNextInputMapping(
				lagertest.NewTestLogger("test"),
				versionsDB,
				fakeJob,
				fakeJobCombination,
			)
		})

		Context("when inputs resolve", func() {
			BeforeEach(func() {
				fakeJob = new(dbfakes.FakeJob)
				fakeJob.NameReturns("some-job")
				fakeJob.ConfigReturns(atc.JobConfig{
					Plan: atc.PlanSequence{
						{Get: "alias", Resource: "a", Version: &atc.VersionConfig{Latest: true}},
						{Get: "b", Version: &atc.VersionConfig{Latest: true}},
					},
				})

				fakeJobCombination = new(dbfakes.FakeJobCombination)
			})

			Context("when transforming the input configs fails", func() {
				BeforeEach(func() {
					fakeTransformer.TransformInputConfigsReturns(nil, disaster)
				})

				It("returns the error", func() {
					Expect(mappingErr).To(Equal(disaster))
				})

				It("transformed the right input configs", func() {
					Expect(fakeTransformer.TransformInputConfigsCallCount()).To(Equal(1))
					actualVersionsDB, actualJobCombination, actualJobInputs := fakeTransformer.TransformInputConfigsArgsForCall(0)
					Expect(actualVersionsDB).To(Equal(versionsDB))
					Expect(actualJobCombination).To(Equal(fakeJobCombination))
					Expect(actualJobInputs).To(ConsistOf(
						atc.JobInput{
							Name:     "alias",
							Resource: "a",
							Version:  &atc.VersionConfig{Latest: true},
						},
						atc.JobInput{
							Name:     "b",
							Resource: "b",
							Version:  &atc.VersionConfig{Latest: true},
						},
					))
				})
			})

			Context("when transforming the input configs succeeds", func() {
				BeforeEach(func() {
					fakeTransformer.TransformInputConfigsReturns(algorithm.InputConfigs{
						{
							Name:             "alias",
							ResourceID:       11,
							Passed:           algorithm.JobSet{},
							JobCombinationID: 1,
						},
						{
							Name:             "b",
							ResourceID:       12,
							Passed:           algorithm.JobSet{},
							JobCombinationID: 1,
						},
					}, nil)
				})

				Context("when saving the independent input mapping fails", func() {
					BeforeEach(func() {
						fakeJobCombination.SaveIndependentInputMappingReturns(disaster)
					})

					It("returns the error", func() {
						Expect(mappingErr).To(Equal(disaster))
					})

					It("saved the right input mapping", func() {
						Expect(fakeJobCombination.SaveIndependentInputMappingCallCount()).To(Equal(1))
						actualMapping := fakeJobCombination.SaveIndependentInputMappingArgsForCall(0)
						Expect(actualMapping).To(Equal(algorithm.InputMapping{
							"alias": algorithm.InputVersion{VersionID: 1, FirstOccurrence: true},
							"b":     algorithm.InputVersion{VersionID: 2, FirstOccurrence: true},
						}))
					})
				})

				Context("when saving the independent input mapping succeeds", func() {
					BeforeEach(func() {
						fakeJobCombination.SaveIndependentInputMappingReturns(nil)
					})

					Context("when saving the next input mapping fails", func() {
						BeforeEach(func() {
							fakeJobCombination.SaveNextInputMappingReturns(disaster)
						})

						It("returns the error", func() {
							Expect(mappingErr).To(Equal(disaster))
						})

						It("saved the right input mapping", func() {
							Expect(fakeJobCombination.SaveIndependentInputMappingCallCount()).To(Equal(1))
							actualMapping := fakeJobCombination.SaveIndependentInputMappingArgsForCall(0)
							Expect(actualMapping).To(Equal(algorithm.InputMapping{
								"alias": algorithm.InputVersion{VersionID: 1, FirstOccurrence: true},
								"b":     algorithm.InputVersion{VersionID: 2, FirstOccurrence: true},
							}))
						})
					})

					Context("when saving the next input mapping succeeds", func() {
						BeforeEach(func() {
							fakeJobCombination.SaveNextInputMappingReturns(nil)
						})

						It("returns the mapping", func() {
							Expect(mappingErr).NotTo(HaveOccurred())
							Expect(inputMapping).To(Equal(algorithm.InputMapping{
								"alias": algorithm.InputVersion{VersionID: 1, FirstOccurrence: true},
								"b":     algorithm.InputVersion{VersionID: 2, FirstOccurrence: true},
							}))
						})

						It("didn't delete the mapping", func() {
							Expect(fakeJobCombination.DeleteNextInputMappingCallCount()).To(BeZero())
						})
					})
				})
			})
		})

		Context("when inputs only resolve individually", func() {
			BeforeEach(func() {
				fakeJob = new(dbfakes.FakeJob)
				fakeJob.NameReturns("some-job")
				fakeJob.ConfigReturns(atc.JobConfig{
					Plan: atc.PlanSequence{
						{Get: "a", Version: &atc.VersionConfig{Latest: true}, Passed: []string{"upstream"}},
						{Get: "b", Version: &atc.VersionConfig{Latest: true}, Passed: []string{"upstream"}},
					},
				})

				fakeJobCombination = new(dbfakes.FakeJobCombination)

				fakeTransformer.TransformInputConfigsReturns(algorithm.InputConfigs{
					{
						Name:             "a",
						ResourceID:       11,
						Passed:           algorithm.JobSet{2: struct{}{}},
						JobCombinationID: 1,
					},
					{
						Name:             "b",
						ResourceID:       12,
						Passed:           algorithm.JobSet{2: struct{}{}},
						JobCombinationID: 1,
					},
				}, nil)
				fakeJobCombination.SaveIndependentInputMappingReturns(nil)
			})

			Context("when deleting the next input mapping fails", func() {
				BeforeEach(func() {
					fakeJobCombination.DeleteNextInputMappingReturns(disaster)
				})

				It("returns the error", func() {
					Expect(mappingErr).To(Equal(disaster))
				})
			})

			Context("when deleting the next input mapping succeeds", func() {
				BeforeEach(func() {
					fakeJobCombination.DeleteNextInputMappingReturns(nil)
				})

				It("saved the right individual input mapping", func() {
					actualMapping := fakeJobCombination.SaveIndependentInputMappingArgsForCall(0)
					Expect(actualMapping).To(Equal(algorithm.InputMapping{
						"a": algorithm.InputVersion{VersionID: 1, FirstOccurrence: true},
						"b": algorithm.InputVersion{VersionID: 2, FirstOccurrence: true},
					}))
				})

				It("deleted the next input mapping", func() {
					Expect(fakeJobCombination.DeleteNextInputMappingCallCount()).To(Equal(1))
					Expect(fakeJobCombination.SaveNextInputMappingCallCount()).To(BeZero())
				})

				It("returns an empty mapping and no error", func() {
					Expect(mappingErr).NotTo(HaveOccurred())
					Expect(inputMapping).To(BeEmpty())
				})
			})
		})

		Context("when some inputs don't resolve", func() {
			BeforeEach(func() {
				fakeJob = new(dbfakes.FakeJob)
				fakeJob.NameReturns("some-job")
				fakeJob.ConfigReturns(atc.JobConfig{
					Plan: atc.PlanSequence{
						{Get: "a", Version: &atc.VersionConfig{Latest: true}},
						{Get: "no-versions", Version: &atc.VersionConfig{Latest: true}},
					},
				})

				fakeJobCombination = new(dbfakes.FakeJobCombination)

				fakeTransformer.TransformInputConfigsReturns(algorithm.InputConfigs{
					{
						Name:             "a",
						ResourceID:       11,
						Passed:           algorithm.JobSet{},
						JobCombinationID: 1,
					},
					{
						Name:             "no-versions",
						ResourceID:       13,
						Passed:           algorithm.JobSet{},
						JobCombinationID: 1,
					},
				}, nil)
				fakeJobCombination.SaveIndependentInputMappingReturns(nil)
				fakeJobCombination.DeleteNextInputMappingReturns(nil)
			})

			It("saved the right individual input mapping", func() {
				actualMapping := fakeJobCombination.SaveIndependentInputMappingArgsForCall(0)
				Expect(actualMapping).To(Equal(algorithm.InputMapping{
					"a": algorithm.InputVersion{VersionID: 1, FirstOccurrence: true},
				}))
			})

			It("deleted the next input mapping", func() {
				Expect(fakeJobCombination.DeleteNextInputMappingCallCount()).To(Equal(1))
				Expect(fakeJobCombination.SaveNextInputMappingCallCount()).To(BeZero())
			})

			It("returns an empty mapping and no error", func() {
				Expect(mappingErr).NotTo(HaveOccurred())
				Expect(inputMapping).To(BeEmpty())
			})
		})

		Context("when a pinned version is missing but the remaining versions resolve", func() {
			BeforeEach(func() {
				fakeJob = new(dbfakes.FakeJob)
				fakeJob.NameReturns("some-job")
				fakeJob.ConfigReturns(atc.JobConfig{
					Plan: atc.PlanSequence{
						{Get: "a", Version: &atc.VersionConfig{Pinned: atc.Version{"doesn't": "exist"}}},
						{Get: "b", Version: &atc.VersionConfig{Latest: true}},
					},
				})

				fakeJobCombination = new(dbfakes.FakeJobCombination)

				fakeTransformer.TransformInputConfigsReturns(algorithm.InputConfigs{
					{
						Name:             "b",
						ResourceID:       12,
						Passed:           algorithm.JobSet{},
						JobCombinationID: 1,
					},
				}, nil)
				fakeJobCombination.SaveIndependentInputMappingReturns(nil)
				fakeJobCombination.DeleteNextInputMappingReturns(nil)
			})

			It("saved the right individual input mapping", func() {
				actualMapping := fakeJobCombination.SaveIndependentInputMappingArgsForCall(0)
				Expect(actualMapping).To(Equal(algorithm.InputMapping{
					"b": algorithm.InputVersion{VersionID: 2, FirstOccurrence: true},
				}))
			})

			It("deleted the next input mapping", func() {
				Expect(fakeJobCombination.DeleteNextInputMappingCallCount()).To(Equal(1))
				Expect(fakeJobCombination.SaveNextInputMappingCallCount()).To(BeZero())
			})

			It("returns an empty mapping and no error", func() {
				Expect(mappingErr).NotTo(HaveOccurred())
				Expect(inputMapping).To(BeEmpty())
			})
		})

		Context("when the job has no inputs", func() {
			BeforeEach(func() {
				fakeJob = new(dbfakes.FakeJob)
				fakeJob.NameReturns("some-job")
				fakeJob.ConfigReturns(atc.JobConfig{
					Plan: atc.PlanSequence{
						{Task: "some-task", TaskConfigPath: "some-task.yml"},
					},
				})

				fakeJobCombination = new(dbfakes.FakeJobCombination)

				fakeTransformer.TransformInputConfigsReturns(algorithm.InputConfigs{}, nil)
				fakeJobCombination.SaveIndependentInputMappingReturns(nil)
				fakeJobCombination.DeleteNextInputMappingReturns(nil)
			})

			It("saved the right individual input mapping", func() {
				actualMapping := fakeJobCombination.SaveIndependentInputMappingArgsForCall(0)
				Expect(actualMapping).To(Equal(algorithm.InputMapping{}))
			})

			It("saved the right next input mapping", func() {
				actualMapping := fakeJobCombination.SaveNextInputMappingArgsForCall(0)
				Expect(actualMapping).To(Equal(algorithm.InputMapping{}))
				Expect(fakeJobCombination.DeleteNextInputMappingCallCount()).To(BeZero())
			})

			It("returns an empty mapping and no error", func() {
				Expect(mappingErr).NotTo(HaveOccurred())
				Expect(inputMapping).To(BeEmpty())
			})
		})
	})
})
