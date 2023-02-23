package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	addCounterRequestDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "counters",
		Subsystem: "http",
		Name:      "add_counter_request_duration",
		Help:      "Add counter request duration",
	}, nil)

	countersNumberGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "counters",
		Subsystem: "http",
		Name:      "counters_number",
		Help:      "",
	}, nil)

	incCounterCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "counters",
		Subsystem: "http",
		Name:      "inc_counter",
		Help:      "Number of requests to increment counter",
	}, nil)

	internalServerErrorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "counters",
		Subsystem: "http",
		Name:      "internal_server_error",
		Help:      "Number of internal server errors",
	}, []string{"method", "uri", "status_code"})
)

func MustRegisterMetrics(registerer prometheus.Registerer) {
	registerer.MustRegister(
		addCounterRequestDurationHistogram,
		countersNumberGauge,
		incCounterCounter,
		internalServerErrorCounter,
	)
}

func withInternalServerErrorCounter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() >= http.StatusInternalServerError {
			defer internalServerErrorCounter.
				With(map[string]string{
					"method":      c.Request.Method,
					"uri":         c.Request.RequestURI,
					"status_code": fmt.Sprintf("%d", c.Writer.Status()),
				}).
				Inc()
		}
	}
}
