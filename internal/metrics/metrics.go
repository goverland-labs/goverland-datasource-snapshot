package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "datasource",
			Name:      "snapshot_sdk_requests",
			Help:      "Time taken to process snapshot requests",
			Buckets:   []float64{.005, .01, .025, .05, .075, .1, .15, .2, .25, .5, 1, 2.5, 5, 10, 15, 30},
		},
		[]string{"client", "method", "error"},
	)

	SnapshotKeyStateGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "datasource",
			Name:      "snapshot_key_state",
			Help:      "Stores the last received number of remaining requests",
		}, []string{"key", "name"},
	)
)

func CollectRequestsMetric(client, method string, err error, start time.Time) {
	RequestsHistogram.
		WithLabelValues(client, method, errLabelValue(err)).
		Observe(time.Since(start).Seconds())
}

func CollectSnapshotKeyState(key, name string, val float64) {
	SnapshotKeyStateGauge.
		WithLabelValues(key, name).
		Set(val)
}

// ErrLabelValue returns string representation of error label value
func errLabelValue(err error) string {
	if err != nil {
		return "true"
	}
	return "false"
}
