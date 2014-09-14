package event_test

import (
	"encoding/json"
	"time"

	"github.com/concourse/turbine/api/builds"
	. "github.com/concourse/turbine/event"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encoding & Decoding messages", func() {
	var event Event

	itEncodesAndDecodesToItself := func() {
		It("encodes and decodes to itself", func() {
			payload, err := json.Marshal(Message{event})
			Ω(err).ShouldNot(HaveOccurred())

			var decodedMsg Message
			err = json.Unmarshal(payload, &decodedMsg)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(decodedMsg.Event).Should(Equal(event))
		})
	}

	Describe("Log", func() {
		BeforeEach(func() {
			event = Log{
				Payload: "some-payload",
				Origin: Origin{
					Type: OriginTypeInput,
					Name: "some-input",
				},
			}
		})

		itEncodesAndDecodesToItself()
	})

	Describe("Status", func() {
		BeforeEach(func() {
			event = Status{
				Status: builds.StatusSucceeded,
			}
		})

		itEncodesAndDecodesToItself()
	})

	Describe("Initialize", func() {
		BeforeEach(func() {
			event = Initialize{
				BuildConfig: builds.Config{
					Image: "some-image",
				},
			}
		})

		itEncodesAndDecodesToItself()
	})

	Describe("Start", func() {
		BeforeEach(func() {
			event = Start{
				Time: time.Now().Unix(),
			}
		})

		itEncodesAndDecodesToItself()
	})

	Describe("Finish", func() {
		BeforeEach(func() {
			event = Finish{
				Time:       time.Now().Unix(),
				ExitStatus: 42,
			}
		})

		itEncodesAndDecodesToItself()
	})

	Describe("Error", func() {
		BeforeEach(func() {
			event = Error{
				Message: "oh no!",
			}
		})

		itEncodesAndDecodesToItself()
	})

	Describe("Input", func() {
		BeforeEach(func() {
			event = Input{
				Input: builds.Input{
					Name: "some-resource",
				},
			}
		})

		itEncodesAndDecodesToItself()
	})

	Describe("Output", func() {
		BeforeEach(func() {
			event = Output{
				Output: builds.Output{
					Name: "some-resource",
				},
			}
		})

		itEncodesAndDecodesToItself()
	})
})
