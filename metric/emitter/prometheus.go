package emitter

import (
	"fmt"
	"net"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/metric"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusEmitter struct {
	buildsStarted     prometheus.Counter
	buildsFinished    prometheus.Counter
	buildsSucceeded   prometheus.Counter
	buildsErrored     prometheus.Counter
	buildsFailed      prometheus.Counter
	buildsAborted     prometheus.Counter
	buildsFinishedVec *prometheus.CounterVec
	buildDurationsVec *prometheus.HistogramVec
}

type PrometheusConfig struct {
	BindIP   string `long:"prometheus-bind-ip" description:"IP to listen on to expose Prometheus metrics."`
	BindPort string `long:"prometheus-bind-port" description:"Port to listen on to expose Prometheus metrics."`
}

func init() {
	metric.RegisterEmitter(&PrometheusConfig{})
}

func (config *PrometheusConfig) Description() string { return "Prometheus" }
func (config *PrometheusConfig) IsConfigured() bool {
	return config.BindPort != "" && config.BindIP != ""
}
func (config *PrometheusConfig) bind() string {
	return fmt.Sprintf("%s:%s", config.BindIP, config.BindPort)
}

func (config *PrometheusConfig) NewEmitter() (metric.Emitter, error) {
	buildsStarted := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "concourse",
		Subsystem: "builds",
		Name:      "started_total",
		Help:      "Total number of Concourse builds started.",
	})
	prometheus.MustRegister(buildsStarted)

	buildsFinished := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "concourse",
		Subsystem: "builds",
		Name:      "finished_total",
		Help:      "Total number of Concourse builds finished.",
	})
	prometheus.MustRegister(buildsFinished)

	buildsSucceeded := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "concourse",
		Subsystem: "builds",
		Name:      "succeeded_total",
		Help:      "Total number of Concourse builds succeeded.",
	})
	prometheus.MustRegister(buildsSucceeded)

	buildsErrored := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "concourse",
		Subsystem: "builds",
		Name:      "errored_total",
		Help:      "Total number of Concourse builds errored.",
	})
	prometheus.MustRegister(buildsErrored)

	buildsFailed := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "concourse",
		Subsystem: "builds",
		Name:      "failed_total",
		Help:      "Total number of Concourse builds failed.",
	})
	prometheus.MustRegister(buildsFailed)

	buildsAborted := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "concourse",
		Subsystem: "builds",
		Name:      "aborted_total",
		Help:      "Total number of Concourse builds aborted.",
	})
	prometheus.MustRegister(buildsAborted)

	buildsFinishedVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "concourse",
			Subsystem: "builds",
			Name:      "finished",
			Help:      "Count of builds finished across various dimensions.",
		},
		[]string{"team", "pipeline", "status"},
	)
	prometheus.MustRegister(buildsFinishedVec)
	buildDurationsVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "concourse",
			Subsystem: "builds",
			Name:      "duration_seconds",
			Help:      "Build time in seconds",
		},
		[]string{"team", "pipeline"},
	)
	prometheus.MustRegister(buildDurationsVec)

	listener, err := net.Listen("tcp", config.bind())
	if err != nil {
		return nil, err
	}

	go http.Serve(listener, promhttp.Handler())

	return &PrometheusEmitter{
		buildsStarted:     buildsStarted,
		buildsFinished:    buildsFinished,
		buildsFinishedVec: buildsFinishedVec,
		buildDurationsVec: buildDurationsVec,
		buildsSucceeded:   buildsSucceeded,
		buildsErrored:     buildsErrored,
		buildsFailed:      buildsFailed,
		buildsAborted:     buildsAborted,
	}, nil
}

// Emit processes incoming metrics.
// In order to provide idiomatic Prometheus metrics, we'll have to convert the various
// Event types (differentiated by the less-than-ideal string Name field) into different
// Prometheus metrics.
func (emitter *PrometheusEmitter) Emit(logger lager.Logger, event metric.Event) {
	switch event.Name {
	case "build started":
		emitter.buildsStarted.Inc()
	case "build finished":
		emitter.buildFinishedMetrics(logger, event)
	default:
		// unless we have a specific metric, we do nothing
	}
}

func (emitter *PrometheusEmitter) buildFinishedMetrics(logger lager.Logger, event metric.Event) {
	// concourse_builds_finished_total
	emitter.buildsFinished.Inc()

	// concourse_builds_finished
	team, exists := event.Attributes["team_name"]
	if !exists {
		logger.Error("failed-to-find-team-name-in-event", fmt.Errorf("expected team_name to exist in event.Attributes"))
	}

	pipeline, exists := event.Attributes["pipeline"]
	if !exists {
		logger.Error("failed-to-find-pipeline-in-event", fmt.Errorf("expected pipeline to exist in event.Attributes"))
	}

	buildStatus, exists := event.Attributes["build_status"]
	if !exists {
		logger.Error("failed-to-find-build_status-in-event", fmt.Errorf("expected build_status to exist in event.Attributes"))
	}
	emitter.buildsFinishedVec.WithLabelValues(team, pipeline, buildStatus).Inc()

	// concourse_builds_(aborted|succeeded|failed|errored)_total
	switch buildStatus {
	case string(db.BuildStatusAborted):
		// concourse_builds_aborted_total
		emitter.buildsAborted.Inc()
	case string(db.BuildStatusSucceeded):
		// concourse_builds_succeeded_total
		emitter.buildsSucceeded.Inc()
	case string(db.BuildStatusFailed):
		// concourse_builds_failed_total
		emitter.buildsFailed.Inc()
	case string(db.BuildStatusErrored):
		// concourse_builds_errored_total
		emitter.buildsErrored.Inc()
	}

	// concourse_builds_duration_seconds
	duration, ok := event.Value.(float64)
	if !ok {
		logger.Error("build-finished-event-value-type-mismatch", fmt.Errorf("expected event.Value to be a float64"))
	}
	// seconds are the standard prometheus base unit for time
	duration = duration / 1000
	emitter.buildDurationsVec.WithLabelValues(team, pipeline).Observe(duration)
}
