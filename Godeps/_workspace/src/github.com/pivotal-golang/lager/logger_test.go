package lager_test

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {
	var logger lager.Logger
	var testSink *lager.TestSink

	var component = "my-component"
	var task = "my-task"
	var action = "my-action"
	var description = "my-description"
	var logData = lager.Data{
		"foo":      "bar",
		"a-number": 7,
	}

	BeforeEach(func() {
		logger = lager.NewLogger(component)
		testSink = lager.NewTestSink()
		logger.RegisterSink(testSink)
	})

	var TestCommonLogFeatures = func(level lager.LogLevel) {
		var log lager.LogFormat

		BeforeEach(func() {
			log = testSink.Logs()[0]
		})

		It("writes a log to the sink", func() {
			Ω(testSink.Logs()).Should(HaveLen(1))
		})

		It("records the source component", func() {
			Ω(log.Source).Should(Equal(component))
		})

		It("outputs a properly-formatted message", func() {
			Ω(log.Message).Should(Equal(fmt.Sprintf("%s.%s.%s", component, task, action)))
		})

		It("has a timestamp", func() {
			expectedTime := float64(time.Now().UnixNano()) / 1e9
			parsedTimestamp, err := strconv.ParseFloat(log.Timestamp, 64)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(parsedTimestamp).Should(BeNumerically("~", expectedTime, 1.0))
		})

		It("data contains the description", func() {
			Ω(log.Data["description"]).Should(Equal(description))
		})

		It("sets the proper output level", func() {
			Ω(log.LogLevel).Should(Equal(level))
		})
	}

	var TestLogData = func() {
		var log lager.LogFormat

		BeforeEach(func() {
			log = testSink.Logs()[0]
		})

		It("data contains custom user data", func() {
			Ω(log.Data["foo"]).Should(Equal("bar"))
			Ω(log.Data["a-number"]).Should(BeNumerically("==", 7))
		})
	}

	Describe("Debug", func() {
		Context("with log data", func() {
			BeforeEach(func() {
				logger.Debug(task, action, description, logData)
			})

			TestCommonLogFeatures(lager.DEBUG)
			TestLogData()
		})

		Context("with no log data", func() {
			BeforeEach(func() {
				logger.Debug(task, action, description)
			})

			TestCommonLogFeatures(lager.DEBUG)
		})
	})

	Describe("Info", func() {
		Context("with log data", func() {
			BeforeEach(func() {
				logger.Info(task, action, description, logData)
			})

			TestCommonLogFeatures(lager.INFO)
			TestLogData()
		})

		Context("with no log data", func() {
			BeforeEach(func() {
				logger.Info(task, action, description)
			})

			TestCommonLogFeatures(lager.INFO)
		})
	})

	Describe("Error", func() {
		var err = errors.New("oh noes!")
		Context("with log data", func() {
			BeforeEach(func() {
				logger.Error(task, action, description, err, logData)
			})

			TestCommonLogFeatures(lager.ERROR)
			TestLogData()

			It("data contains error message", func() {
				Ω(testSink.Logs()[0].Data["error"]).Should(Equal(err.Error()))
			})
		})

		Context("with no log data", func() {
			BeforeEach(func() {
				logger.Error(task, action, description, err)
			})

			TestCommonLogFeatures(lager.ERROR)

			It("data contains error message", func() {
				Ω(testSink.Logs()[0].Data["error"]).Should(Equal(err.Error()))
			})
		})

		Context("with no error", func() {
			BeforeEach(func() {
				logger.Error(task, action, description, nil)
			})

			TestCommonLogFeatures(lager.ERROR)

			It("does not contain the error message", func() {
				Ω(testSink.Logs()[0].Data).ShouldNot(HaveKey("error"))
			})
		})
	})

	Describe("Fatal", func() {
		var err = errors.New("oh noes!")
		var fatalErr interface{}

		Context("with log data", func() {
			BeforeEach(func() {
				defer func() {
					fatalErr = recover()
				}()

				logger.Fatal(task, action, description, err, logData)
			})

			TestCommonLogFeatures(lager.FATAL)
			TestLogData()

			It("data contains error message", func() {
				Ω(testSink.Logs()[0].Data["error"]).Should(Equal(err.Error()))
			})

			It("data contains stack trace", func() {
				Ω(testSink.Logs()[0].Data["trace"]).ShouldNot(BeEmpty())
			})

			It("panics with the provided error", func() {
				Ω(fatalErr).Should(Equal(err))
			})
		})

		Context("with no log data", func() {
			BeforeEach(func() {
				defer func() {
					fatalErr = recover()
				}()

				logger.Fatal(task, action, description, err)
			})

			TestCommonLogFeatures(lager.FATAL)

			It("data contains error message", func() {
				Ω(testSink.Logs()[0].Data["error"]).Should(Equal(err.Error()))
			})

			It("data contains stack trace", func() {
				Ω(testSink.Logs()[0].Data["trace"]).ShouldNot(BeEmpty())
			})

			It("panics with the provided error", func() {
				Ω(fatalErr).Should(Equal(err))
			})
		})

		Context("with no error", func() {
			BeforeEach(func() {
				defer func() {
					fatalErr = recover()
				}()

				logger.Fatal(task, action, description, nil)
			})

			TestCommonLogFeatures(lager.FATAL)

			It("does not contain the error message", func() {
				Ω(testSink.Logs()[0].Data).ShouldNot(HaveKey("error"))
			})

			It("data contains stack trace", func() {
				Ω(testSink.Logs()[0].Data["trace"]).ShouldNot(BeEmpty())
			})

			It("panics with the provided error (i.e. nil)", func() {
				Ω(fatalErr).Should(BeNil())
			})
		})
	})
})
