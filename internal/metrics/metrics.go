package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestCounter *prometheus.CounterVec

func init() {
	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_counter",
		Help: "A counter of total requests by status code",
	}, []string{"code", "method"})
}

func RequestInc(statusCode int, method string) {
	requestCounter.With(prometheus.Labels{"code": strconv.Itoa(statusCode), "method": method}).Inc()
}
