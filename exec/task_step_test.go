package exec_test

import (
	"context"
	"errors"

	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/dbfakes"
	"github.com/concourse/atc/exec"
	"github.com/concourse/atc/exec/execfakes"
	"github.com/concourse/atc/runtime"
	"github.com/concourse/atc/runtime/runtimefakes"
	"github.com/concourse/atc/worker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("TaskStep", func() {
	var (
		ctx    context.Context
		cancel func()

		fakeOrchestrator           *runtimefakes.FakeOrchestrator
		fakeDBResourceCacheFactory *dbfakes.FakeResourceCacheFactory

		stdoutBuf *gbytes.Buffer
		stderrBuf *gbytes.Buffer

		imageArtifactName string
		containerMetadata db.ContainerMetadata

		fakeDelegate *execfakes.FakeTaskDelegate

		privileged    exec.Privileged
		tags          []string
		teamID        int
		buildID       int
		planID        atc.PlanID
		jobID         int
		configSource  *execfakes.FakeTaskConfigSource
		resourceTypes creds.VersionedResourceTypes
		inputMapping  map[string]string
		outputMapping map[string]string
		variables     creds.Variables

		repo  *worker.ArtifactRepository
		state *execfakes.FakeRunState

		taskStep exec.Step

		stepErr    error
		runFunc    func()
		taskResult runtime.TaskResult
		mounts     []worker.VolumeMount
	)

	var result = make(chan runtime.TaskResult)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		fakeOrchestrator = new(runtimefakes.FakeOrchestrator)
		fakeDBResourceCacheFactory = new(dbfakes.FakeResourceCacheFactory)

		stdoutBuf = gbytes.NewBuffer()
		stderrBuf = gbytes.NewBuffer()

		fakeDelegate = new(execfakes.FakeTaskDelegate)
		fakeDelegate.StdoutReturns(stdoutBuf)
		fakeDelegate.StderrReturns(stderrBuf)

		privileged = false
		tags = []string{"step", "tags"}
		teamID = 123
		planID = atc.PlanID(42)
		buildID = 1234
		jobID = 12345
		configSource = new(execfakes.FakeTaskConfigSource)

		repo = worker.NewArtifactRepository()
		state = new(execfakes.FakeRunState)
		state.ArtifactsReturns(repo)

		resourceTypes = creds.NewVersionedResourceTypes(variables, atc.VersionedResourceTypes{
			{
				ResourceType: atc.ResourceType{
					Name:   "custom-resource",
					Type:   "custom-type",
					Source: atc.Source{"some-custom": "source"},
					Params: atc.Params{"some-custom": "param"},
				},
				Version: atc.Version{"some-custom": "version"},
			},
		})

		inputMapping = nil
		outputMapping = nil
		imageArtifactName = ""

		variables = template.StaticVariables{
			"source-param": "super-secret-source",
			"task-param":   "super-secret-param",
		}

		containerMetadata = db.ContainerMetadata{
			Type:     db.ContainerTypeTask,
			StepName: "some-step",
		}

		fakeOrchestrator.RunTaskReturns(result, mounts, nil)

		stepErr = nil

		runFunc = func() {
			// this guards against ny boo-boos when running the task and not sending
			// to the result channel or cancelling or otherwise closing the context
			stepExited := make(chan bool)
			go func() {
				stepErr = taskStep.Run(ctx, state)
				stepExited <- true
			}()

			Eventually(stepExited).Should(Receive())
		}
	})

	JustBeforeEach(func() {
		taskStep = exec.NewTaskStep(
			privileged,
			configSource,
			tags,
			inputMapping,
			outputMapping,
			"some-artifact-root",
			imageArtifactName,
			fakeDelegate,
			fakeOrchestrator,
			teamID,
			buildID,
			jobID,
			"some-task",
			planID,
			containerMetadata,
			resourceTypes,
			variables,
		)
		runFunc()
	})

	Context("when getting the config works", func() {
		var fetchedConfig atc.TaskConfig

		BeforeEach(func() {
			fetchedConfig = atc.TaskConfig{
				Platform: "some-platform",
				ImageResource: &atc.ImageResource{
					Type:    "docker",
					Source:  atc.Source{"some": "((source-param))"},
					Params:  &atc.Params{"some": "params"},
					Version: &atc.Version{"some": "version"},
				},
				Params: map[string]string{
					"SECURE": "((task-param))",
				},
				Run: atc.TaskRunConfig{
					Path: "ls",
					Args: []string{"some", "args"},
				},
			}

			configSource.FetchConfigReturns(fetchedConfig, nil)
		})

		Describe("before orchestration of the task", func() {
			BeforeEach(func() {
				fakeDelegate.InitializingStub = func(lager.Logger, atc.TaskConfig) {
					defer GinkgoRecover()
					Expect(fakeOrchestrator.RunTaskCallCount()).To(BeZero())
					cancel()
				}
			})

			It("invokes the delegate's Initializing callback", func() {
				Expect(fakeDelegate.InitializingCallCount()).To(Equal(1))
			})
		})

		Context("when the orchestrator is able to run the Task", func() {
			BeforeEach(func() {
				runFunc = func() {

					stepExited := make(chan bool)
					go func() {
						stepErr = taskStep.Run(ctx, state)
						stepExited <- true
					}()

					go func() {
						result <- taskResult
					}()

					Eventually(stepExited).Should(Receive())
				}
			})

			Context("when the process errors", func() {
				var procesErr = errors.New("failure")
				BeforeEach(func() {
					taskResult.Err = procesErr
				})

				It("returns errors from the task result ", func() {
					Expect(stepErr).To(Equal(procesErr))
				})
			})
		})

	})

	Context("when task orchestration fails", func() {
		disaster := errors.New("nope")

		BeforeEach(func() {
			fakeOrchestrator.RunTaskReturns(nil, nil, disaster)
		})

		It("returns the error", func() {
			Expect(stepErr).To(Equal(disaster))
		})
	})

	Context("when getting the config fails", func() {
		disaster := errors.New("nope")

		BeforeEach(func() {
			configSource.FetchConfigReturns(atc.TaskConfig{}, disaster)
		})

		It("returns the error", func() {
			Expect(stepErr).To(Equal(disaster))
		})
	})
})
