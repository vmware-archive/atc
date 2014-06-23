package lager

import (
	"fmt"
	"runtime"
	"time"
)

const STACK_TRACE_BUFFER_SIZE = 1024 * 100

type Logger interface {
	RegisterSink(Sink)
	Debug(task, action, description string, data ...Data)
	Info(task, action, description string, data ...Data)
	Error(task, action, description string, err error, data ...Data)
	Fatal(task, action, description string, err error, data ...Data)
}

type logger struct {
	component string
	sinks     []Sink
}

func NewLogger(component string) Logger {
	return &logger{
		component: component,
		sinks:     []Sink{},
	}
}

func (l *logger) RegisterSink(sink Sink) {
	l.sinks = append(l.sinks, sink)
}

func (l *logger) Debug(task, action, description string, data ...Data) {
	logData := Data{}
	if len(data) > 0 {
		logData = data[0]
	}

	logData["description"] = description

	log := LogFormat{
		Timestamp: currentTimestamp(),
		Source:    l.component,
		Message:   fmt.Sprintf("%s.%s.%s", l.component, task, action),
		LogLevel:  DEBUG,
		Data:      logData,
	}

	for _, sink := range l.sinks {
		sink.Log(DEBUG, log.ToJSON())
	}
}

func (l *logger) Info(task, action, description string, data ...Data) {
	logData := Data{}
	if len(data) > 0 {
		logData = data[0]
	}

	logData["description"] = description

	log := LogFormat{
		Timestamp: currentTimestamp(),
		Source:    l.component,
		Message:   fmt.Sprintf("%s.%s.%s", l.component, task, action),
		LogLevel:  INFO,
		Data:      logData,
	}

	for _, sink := range l.sinks {
		sink.Log(INFO, log.ToJSON())
	}
}

func (l *logger) Error(task, action, description string, err error, data ...Data) {
	logData := Data{}
	if len(data) > 0 {
		logData = data[0]
	}

	logData["description"] = description
	if err != nil {
		logData["error"] = err.Error()
	}

	log := LogFormat{
		Timestamp: currentTimestamp(),
		Source:    l.component,
		Message:   fmt.Sprintf("%s.%s.%s", l.component, task, action),
		LogLevel:  ERROR,
		Data:      logData,
	}

	for _, sink := range l.sinks {
		sink.Log(ERROR, log.ToJSON())
	}
}

func (l *logger) Fatal(task, action, description string, err error, data ...Data) {
	logData := Data{}
	if len(data) > 0 {
		logData = data[0]
	}

	stackTrace := make([]byte, STACK_TRACE_BUFFER_SIZE)
	stackSize := runtime.Stack(stackTrace, false)
	stackTrace = stackTrace[:stackSize]

	logData["description"] = description
	if err != nil {
		logData["error"] = err.Error()
	}
	logData["trace"] = string(stackTrace)

	log := LogFormat{
		Timestamp: currentTimestamp(),
		Source:    l.component,
		Message:   fmt.Sprintf("%s.%s.%s", l.component, task, action),
		LogLevel:  FATAL,
		Data:      logData,
	}

	for _, sink := range l.sinks {
		sink.Log(FATAL, log.ToJSON())
	}

	panic(err)
}

func currentTimestamp() string {
	return fmt.Sprintf("%.9f", float64(time.Now().UnixNano())/1e9)
}
