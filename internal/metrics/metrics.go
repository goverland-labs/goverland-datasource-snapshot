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
			Buckets:   []float64{20, 50, 100, 500},
		},
		[]string{"client", "method", "error"},
	)
)

func CollectRequestsMetric(client, method string, err error, start time.Time) {
	RequestsHistogram.
		WithLabelValues(client, method, errLabelValue(err)).
		Observe(time.Since(start).Seconds())
}

// ErrLabelValue returns string representation of error label value
func errLabelValue(err error) string {
	if err != nil {
		return "true"
	}
	return "false"
}
